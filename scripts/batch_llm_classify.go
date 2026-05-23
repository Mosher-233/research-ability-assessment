package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Mosher-233/research-ability-assessment/internal/agent"
	"github.com/Mosher-233/research-ability-assessment/internal/config"
	"github.com/Mosher-233/research-ability-assessment/internal/llm"
	"github.com/Mosher-233/research-ability-assessment/internal/service"
	"github.com/Mosher-233/research-ability-assessment/pkg/extractor"
)

// Level mapping from filename prefix to expected level.
var prefixToLevel = map[string]int{
	"A": 4,
	"B": 3,
	"C": 2,
	"D": 1,
	"E": 1, // 答非所问 → should be flagged
	"F": 1, // 内容空洞 → should be flagged
	"G": 1, // 格式混乱 → should be flagged
	"H": 1, // 异常输入 → should be flagged
}

var prefixToLabel = map[string]string{
	"A": "优秀",
	"B": "良好",
	"C": "合格",
	"D": "不合格",
	"E": "答非所问",
	"F": "内容空洞",
	"G": "格式混乱",
	"H": "异常输入",
}

var anomalousPrefixes = map[string]bool{
	"E": true, "F": true, "G": true, "H": true,
}

type BatchResult struct {
	Filename      string
	SourceType    string // "docx" or "pdf"
	Prefix        string // A-H
	ExpectedLevel int
	IsAnomalous   bool
	CharCount     int
	KBMName       string
	Level         int
	Credibility   float64
	Rationale     string
	Error         string
	Duration      time.Duration
}

func main() {
	configPath := flag.String("config", "configs/config.dev.yaml", "path to config YAML")
	sample := flag.Int("sample", 0, "limit to N files (0 = all)")
	workers := flag.Int("workers", 5, "concurrent LLM workers")
	flag.Parse()

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("  LLM 批量证据分类与等级评定测试")
	fmt.Println("=" + strings.Repeat("=", 60))

	// Load config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("LLM Provider: %s | Model: %s | BaseURL: %s\n", cfg.LLM.Provider, cfg.LLM.Model, cfg.LLM.BaseURL)
	fmt.Printf("并发Worker: %d\n\n", *workers)

	// Create LLM client and agent
	llmClient := llm.NewClient(&cfg.LLM)
	evidenceService := service.NewEvidenceService(nil, llmClient)
	evidenceAgent := agent.NewEvidenceAgentWithLLM(evidenceService, llmClient)

	// Collect all files
	files := collectFiles()
	fmt.Printf("收集到 %d 个测试文件\n", len(files))

	if *sample > 0 && *sample < len(files) {
		files = files[:*sample]
		fmt.Printf("采样模式: 仅处理前 %d 个文件\n", *sample)
	}

	// Process files concurrently
	results := processFiles(files, evidenceAgent, *workers)

	// Analyze and print report
	printReport(results, files)
}

type fileInfo struct {
	Path       string
	SourceType string // "docx" or "pdf"
	Prefix     string
}

func collectFiles() []fileInfo {
	var files []fileInfo

	// DOCX files
	docxMatches, _ := filepath.Glob("testdata/words/*.docx")
	for _, p := range docxMatches {
		base := filepath.Base(p)
		prefix := extractPrefix(base)
		if prefix != "" {
			files = append(files, fileInfo{Path: p, SourceType: "docx", Prefix: prefix})
		}
	}

	// PDF files
	pdfMatches, _ := filepath.Glob("testdata/pdfs/*.pdf")
	for _, p := range pdfMatches {
		base := filepath.Base(p)
		prefix := extractPrefix(base)
		if prefix != "" {
			files = append(files, fileInfo{Path: p, SourceType: "pdf", Prefix: prefix})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Prefix < files[j].Prefix ||
			(files[i].Prefix == files[j].Prefix && files[i].Path < files[j].Path)
	})

	return files
}

func extractPrefix(filename string) string {
	if idx := strings.Index(filename, "_"); idx > 0 {
		p := filename[:idx]
		if len(p) == 1 && p[0] >= 'A' && p[0] <= 'H' {
			return p
		}
	}
	return ""
}

func processFiles(files []fileInfo, evidenceAgent *agent.EvidenceAgent, workers int) []BatchResult {
	extractorChain := extractor.NewExtractorChain()
	results := make([]BatchResult, len(files))

	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)
	var processed atomic.Int64
	total := len(files)

	for i, f := range files {
		wg.Add(1)
		go func(idx int, fi fileInfo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := BatchResult{
				Filename:   filepath.Base(fi.Path),
				SourceType: fi.SourceType,
				Prefix:     fi.Prefix,
				ExpectedLevel: prefixToLevel[fi.Prefix],
				IsAnomalous:   anomalousPrefixes[fi.Prefix],
			}

			start := time.Now()

			// Extract text
			extractResult := extractorChain.Extract(fi.Path)
			if extractResult == nil || strings.TrimSpace(extractResult.Content) == "" {
				result.Error = "extraction failed or empty content"
				result.Duration = time.Since(start)
				results[idx] = result
				n := processed.Add(1)
				fmt.Printf("\r[%3d/%3d] %s: 提取失败", n, total, result.Filename)
				return
			}

			result.CharCount = len([]rune(extractResult.Content))

			// Run LLM classification + assessment
			kbmInfo := evidenceAgent.ExtractKBMInfo(extractResult.Content)
			result.KBMName = kbmInfo.KBMName
			result.Level = kbmInfo.Level
			result.Credibility = kbmInfo.Credibility
			result.Rationale = kbmInfo.Rationale
			result.Duration = time.Since(start)

			results[idx] = result
			n := processed.Add(1)
			fmt.Printf("\r[%3d/%3d] %s → %s L%d (%.0f%%)", n, total, result.Filename, result.KBMName, result.Level, result.Credibility*100)
		}(i, f)
	}

	wg.Wait()
	fmt.Println()
	return results
}

func printReport(results []BatchResult, files []fileInfo) {
	success := 0
	failed := 0
	for _, r := range results {
		if r.Error != "" {
			failed++
		} else {
			success++
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  测试结果汇总")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("总文件数: %d | 成功: %d | 失败: %d | 成功率: %.1f%%\n",
		len(results), success, failed, float64(success)/float64(len(results))*100)

	// 1. KBM distribution
	fmt.Println("\n## KBM 分类分布")
	kbmCount := make(map[string]int)
	for _, r := range results {
		if r.Error == "" {
			kbmCount[r.KBMName]++
		}
	}
	type kv struct{ k string; v int }
	var kbmList []kv
	for k, v := range kbmCount {
		kbmList = append(kbmList, kv{k, v})
	}
	sort.Slice(kbmList, func(i, j int) bool { return kbmList[i].v > kbmList[j].v })
	fmt.Println("| KBM类别 | 数量 | 占比 |")
	fmt.Println("|---------|------|------|")
	for _, item := range kbmList {
		fmt.Printf("| %s | %d | %.1f%% |\n", item.k, item.v, float64(item.v)/float64(success)*100)
	}

	// 2. Level assessment by prefix
	fmt.Println("\n## 按质量等级分组评定")
	fmt.Println("| 等级 | 文件数 | LLM平均Level | Level±0一致率 | Level±1以内率 | 平均可信度 |")
	fmt.Println("|------|--------|-------------|--------------|--------------|-----------|")
	prefixes := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, p := range prefixes {
		var group []BatchResult
		for _, r := range results {
			if r.Prefix == p && r.Error == "" {
				group = append(group, r)
			}
		}
		if len(group) == 0 {
			continue
		}
		var sumLevel float64
		var sumCred float64
		exactMatch := 0
		withinOne := 0
		expLevel := prefixToLevel[p]
		for _, r := range group {
			sumLevel += float64(r.Level)
			sumCred += r.Credibility
			if r.Level == expLevel {
				exactMatch++
			}
			diff := r.Level - expLevel
			if diff < 0 {
				diff = -diff
			}
			if diff <= 1 {
				withinOne++
			}
		}
		n := len(group)
		fmt.Printf("| %s (%s) | %d | %.2f | %d/%d (%.0f%%) | %d/%d (%.0f%%) | %.0f%% |\n",
			p, prefixToLabel[p], n,
			sumLevel/float64(n),
			exactMatch, n, float64(exactMatch)/float64(n)*100,
			withinOne, n, float64(withinOne)/float64(n)*100,
			sumCred/float64(n)*100,
		)
	}

	// 3. Global accuracy metrics
	fmt.Println("\n## 综合准确度指标")
	exactMatch := 0
	withinOne := 0
	overestimated := 0
	underestimated := 0
	var sumAbsDiff float64
	validResults := 0
	for _, r := range results {
		if r.Error != "" || r.IsAnomalous {
			continue // skip anomalous for level accuracy (they have no "correct" level)
		}
		validResults++
		diff := r.Level - r.ExpectedLevel
		if diff < 0 {
			diff = -diff
		}
		sumAbsDiff += float64(diff)
		if r.Level == r.ExpectedLevel {
			exactMatch++
		}
		if diff <= 1 {
			withinOne++
		}
		if r.Level > r.ExpectedLevel {
			overestimated++
		} else if r.Level < r.ExpectedLevel {
			underestimated++
		}
	}

	if validResults > 0 {
		fmt.Println("| 指标 | 数值 |")
		fmt.Println("|------|------|")
		fmt.Printf("| 有效样本 (非异常) | %d |\n", validResults)
		fmt.Printf("| 等级完全一致 | %d/%d (%.1f%%) |\n", exactMatch, validResults, float64(exactMatch)/float64(validResults)*100)
		fmt.Printf("| 等级±1以内 | %d/%d (%.1f%%) |\n", withinOne, validResults, float64(withinOne)/float64(validResults)*100)
		fmt.Printf("| 高估数量 | %d (%.1f%%) |\n", overestimated, float64(overestimated)/float64(validResults)*100)
		fmt.Printf("| 低估数量 | %d (%.1f%%) |\n", underestimated, float64(underestimated)/float64(validResults)*100)
		fmt.Printf("| 平均绝对偏差 (MAE) | %.2f |\n", sumAbsDiff/float64(validResults))
	}

	// 4. Anomalous content detection
	fmt.Println("\n## 异常内容检测 (E/F/G/H)")
	anomalyDetected := 0
	anomalyTotal := 0
	for _, r := range results {
		if r.IsAnomalous && r.Error == "" {
			anomalyTotal++
			// Consider "detected" if level is L1 with low credibility
			if r.Level == 1 && r.Credibility < 0.5 {
				anomalyDetected++
			} else if r.Level == 1 {
				anomalyDetected++
			}
		}
	}
	if anomalyTotal > 0 {
		fmt.Printf("| 异常文件总数 | %d |\n", anomalyTotal)
		fmt.Printf("| 正确评定为L1 | %d (%.0f%%) |\n", anomalyDetected, float64(anomalyDetected)/float64(anomalyTotal)*100)
	}

	// 5. Duration stats
	fmt.Println("\n## 性能统计")
	var sumDur time.Duration
	var minDur, maxDur time.Duration = 1<<63 - 1, 0
	for _, r := range results {
		if r.Error == "" {
			sumDur += r.Duration
			if r.Duration < minDur {
				minDur = r.Duration
			}
			if r.Duration > maxDur {
				maxDur = r.Duration
			}
		}
	}
	if success > 0 {
		fmt.Printf("| 总耗时 | %v |\n", sumDur)
		fmt.Printf("| 平均响应时间 | %v |\n", sumDur/time.Duration(success))
		fmt.Printf("| 最短 | %v |\n", minDur)
		fmt.Printf("| 最长 | %v |\n", maxDur)
	}

	// 6. Failed items
	if failed > 0 {
		fmt.Println("\n## 失败项")
		for _, r := range results {
			if r.Error != "" {
				fmt.Printf("- **%s**: %s\n", r.Filename, r.Error)
			}
		}
	}

	// 7. Per-file detail for anomalous items
	fmt.Println("\n## 异常文件LLM处理详情")
	for _, r := range results {
		if r.IsAnomalous && r.Error == "" {
			fmt.Printf("- **%s** [%s]: KBM=%s, Level=L%d, Cred=%.0f%% | %s\n",
				r.Filename, prefixToLabel[r.Prefix],
				r.KBMName, r.Level, r.Credibility*100, truncateStr(r.Rationale, 150))
		}
	}
}

func truncateStr(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
