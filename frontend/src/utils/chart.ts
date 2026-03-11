// 生成雷达图数据
export const generateRadarData = (dimensionScores: Record<string, { score: number }>) => {
  const dimensions = Object.keys(dimensionScores)
  const scores = dimensions.map(dim => dimensionScores[dim].score)
  
  return {
    radar: {
      indicator: dimensions.map(dim => ({
        name: dim,
        max: 1
      }))
    },
    series: [{
      data: [{
        value: scores,
        name: '研究能力'
      }]
    }]
  }
}

// 生成柱状图数据
export const generateBarData = (results: Array<{ student: { name: string }; overall_score: number }>) => {
  const names = results.map(r => r.student.name)
  const scores = results.map(r => r.overall_score)
  
  return {
    xAxis: {
      type: 'category',
      data: names
    },
    yAxis: {
      type: 'value',
      max: 1
    },
    series: [{
      data: scores,
      type: 'bar'
    }]
  }
}

// 生成饼图数据
export const generatePieData = (dimensionScores: Record<string, { score: number }>) => {
  return Object.entries(dimensionScores).map(([name, score]) => ({
    name,
    value: score.score
  }))
}