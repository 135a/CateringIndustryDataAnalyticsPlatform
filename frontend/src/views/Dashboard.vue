<template>
  <div class="dashboard-container">
    <!-- 头部导航栏：项目大标题与用户信息 -->
    <header class="header glass-panel">
      <h1>餐饮大数据平台实时监控中心</h1>
      <div class="user-info">
        <span class="greeting">欢迎回来，{{ username }}</span>
        <button class="logout-btn" @click="handleLogout">退出登录</button>
      </div>
    </header>

    <!-- 核心业务指标卡片 (Overview) -->
    <div class="metrics-grid">
      <div class="metric-card glass-panel" v-for="(val, key) in metricsLabels" :key="key">
        <div class="metric-title">{{ val.title }}</div>
        <div class="metric-value">
          {{ formatNumber(overview[key]) }}
          <span class="metric-unit">{{ val.unit }}</span>
        </div>
      </div>
    </div>

    <!-- 核心图表展示区 -->
    <div class="charts-grid">
      <!-- 饼图/玫瑰图 -->
      <div class="chart-box glass-panel">
        <div class="chart-title">热门品类占比</div>
        <div ref="pieChartRef" class="chart"></div>
      </div>
      
      <!-- 柱状图 + 折线图混合 -->
      <div class="chart-box glass-panel span-2">
        <div class="chart-title">杭州各行政区商户与消费双轴对比</div>
        <div ref="barChartRef" class="chart"></div>
      </div>

      <!-- 大范围散点图 -->
      <div class="chart-box glass-panel span-full">
        <div class="chart-title">价格与评分性价比分布 (发掘高分低价宝藏店)</div>
        <div ref="scatterChartRef" class="chart" style="height: 400px;"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import request from '../utils/request'

const router = useRouter()
// 尝试从浏览器缓存读取用户名，默认为 Admin
const username = ref(localStorage.getItem('username') || 'Admin')

// 安全退出逻辑
const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}

// ================= 指标卡数据模型 =================
const metricsLabels = {
  total_shops: { title: '收录商户总数', unit: '家' },
  avg_price: { title: '全市人均消费', unit: '元' },
  total_reviews: { title: '累计评价总量', unit: '条' },
  avg_rating: { title: '全网平均星级', unit: '星' }
}
const overview = ref({
  total_shops: 0,
  avg_price: 0,
  total_reviews: 0,
  avg_rating: 0
})

// 数字格式化工具 (如：1234567 -> 1,234,567)
const formatNumber = (num) => {
  if (!num) return 0
  return Number(num).toLocaleString(undefined, { maximumFractionDigits: 1 })
}

// ================= ECharts 挂载与渲染逻辑 =================
const pieChartRef = ref(null)
const barChartRef = ref(null)
const scatterChartRef = ref(null)

let pieChart, barChart, scatterChart
const textStyle = { color: '#94a3b8' } // 统一暗色系文字色

const loadData = async () => {
  try {
    // 1. 并发请求四个核心接口，提升数据加载速度
    const [overRes, pieRes, barRes, scatterRes] = await Promise.all([
      request.get('/statistics/overview'),
      request.get('/statistics/category-pie'),
      request.get('/statistics/district-bar'),
      request.get('/statistics/price-rating-scatter')
    ])

    // 2. 绑定指标卡数据
    overview.value = overRes

    // 3. 渲染炫酷的南丁格尔玫瑰图
    renderPieChart(pieRes)

    // 4. 渲染柱状折线双轴图
    renderBarChart(barRes)

    // 5. 渲染绿宝石散点图
    renderScatterChart(scatterRes)

  } catch (error) {
    console.error("加载大屏数据失败", error)
  }
}

// 图表渲染器实现细节...
const renderPieChart = (data) => {
  pieChart = echarts.init(pieChartRef.value)
  pieChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b} : {c}家 ({d}%)' },
    series: [
      {
        name: '品类占比',
        type: 'pie',
        radius: [20, 100], // 空心环状效果
        center: ['50%', '50%'],
        roseType: 'area',  // 开启玫瑰图模式
        itemStyle: { borderRadius: 6 },
        label: { color: '#fff' },
        data: data
      }
    ]
  })
}

const renderBarChart = (data) => {
  barChart = echarts.init(barChartRef.value)
  barChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { textStyle, top: 0 },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.xAxis,
      axisLabel: textStyle
    },
    yAxis: [
      { 
        type: 'value', name: '商户数', 
        axisLabel: textStyle, nameTextStyle: textStyle, 
        splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } } 
      },
      { 
        type: 'value', name: '均价(元)', 
        axisLabel: textStyle, nameTextStyle: textStyle, 
        splitLine: { show: false } 
      }
    ],
    series: [
      {
        name: '商户总数',
        type: 'bar',
        data: data.shop_counts,
        itemStyle: {
          // 渐变蓝柱体
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#0ea5e9' },
            { offset: 1, color: '#3b82f6' }
          ]),
          borderRadius: [4, 4, 0, 0]
        }
      },
      {
        name: '人均消费',
        type: 'line',
        yAxisIndex: 1,
        data: data.avg_prices,
        smooth: true,
        lineStyle: { width: 3, color: '#f59e0b' },
        itemStyle: { color: '#f59e0b' }
      }
    ]
  })
}

const renderScatterChart = (data) => {
  scatterChart = echarts.init(scatterChartRef.value)
  scatterChart.setOption({
    tooltip: {
      backgroundColor: 'rgba(20, 26, 40, 0.9)',
      borderColor: '#3b82f6',
      textStyle: { color: '#fff' },
      formatter: function (params) {
        // [价格, 评分, 店名, 分类]
        return `<div style="padding: 4px;">
                  <b style="color: #34d399; font-size: 16px;">${params.value[2]}</b> 
                  <span style="color: #94a3b8; font-size: 12px;">(${params.value[3]})</span><hr style="border:0;border-top:1px solid rgba(255,255,255,0.1);margin:6px 0;"/>
                  💰 人均: ¥${params.value[0]}<br/>
                  ⭐ 评分: ${params.value[1]} 星
                </div>`
      }
    },
    grid: { left: '3%', right: '5%', bottom: '3%', containLabel: true },
    xAxis: {
      name: '人均价格 (元)',
      type: 'value',
      axisLabel: textStyle, nameTextStyle: textStyle,
      splitLine: { show: false }
    },
    yAxis: {
      name: '大众评分',
      type: 'value',
      min: 2, max: 5, // 评分一般在 2 到 5 之间
      axisLabel: textStyle, nameTextStyle: textStyle,
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
    },
    series: [
      {
        symbolSize: 8,
        data: data,
        type: 'scatter',
        itemStyle: {
          color: 'rgba(16, 185, 129, 0.7)',
          shadowBlur: 10,
          shadowColor: 'rgba(16, 185, 129, 0.4)'
        }
      }
    ]
  })
}

// 绑定与解绑自适应窗口改变事件
onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  pieChart?.dispose()
  barChart?.dispose()
  scatterChart?.dispose()
})

const handleResize = () => {
  pieChart?.resize()
  barChart?.resize()
  scatterChart?.resize()
}
</script>

<style scoped>
.dashboard-container {
  padding: 1.5rem;
  max-width: 1600px;
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 顶部导航 */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.2rem 2rem;
}

.header h1 {
  font-size: 1.6rem;
  font-weight: 700;
  letter-spacing: 2px;
  background: linear-gradient(to right, #60a5fa, #34d399);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.greeting {
  color: var(--text-secondary);
  font-weight: 500;
}

.logout-btn {
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.4);
  color: #fca5a5;
  padding: 0.4rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: var(--transition-fast);
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: #ef4444;
}

/* 指标卡片阵列 */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.metric-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border-top: 3px solid var(--brand-primary);
  transition: transform var(--transition-fast);
}

.metric-card:hover {
  transform: translateY(-5px);
}

.metric-title {
  color: var(--text-muted);
  font-size: 0.95rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.metric-value {
  font-size: 2.2rem;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: baseline;
  gap: 0.3rem;
  /* 科技感数字阴影 */
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.2);
}

.metric-unit {
  font-size: 1rem;
  color: var(--text-secondary);
  font-weight: 400;
}

/* 图表矩阵 */
.charts-grid {
  display: grid;
  grid-template-columns: 1fr 2.5fr; /* 饼图占1份，双轴图占2.5份 */
  gap: 1.5rem;
}

.chart-box {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
}

.chart-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 1rem;
  padding-left: 0.8rem;
  border-left: 4px solid var(--brand-accent);
}

.chart {
  flex: 1;
  min-height: 320px;
}

.span-2 {
  grid-column: span 1;
}

.span-full {
  grid-column: 1 / -1;
}

/* 响应式降级 */
@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
