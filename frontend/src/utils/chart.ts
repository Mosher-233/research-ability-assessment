// 生成雷达图数据
export const generateRadarData = (dimensionScores: Record<string, { score: number }>) => {
  const dimensions = Object.keys(dimensionScores)
  const scores = dimensions.map(dim => {
    const score = dimensionScores[dim].score
    return score <= 1 ? score * 100 : score
  })
  
  return {
    radar: {
      indicator: dimensions.map(dim => ({
        name: dim,
        max: 100
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
  const scores = results.map(r => {
    const score = r.overall_score
    return score <= 1 ? score * 100 : score
  })
  
  return {
    xAxis: {
      type: 'category',
      data: names
    },
    yAxis: {
      type: 'value',
      max: 100
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