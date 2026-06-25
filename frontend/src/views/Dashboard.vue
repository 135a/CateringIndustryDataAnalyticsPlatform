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
      <div 
        class="metric-card glass-panel" 
        v-for="(val, key) in metricsLabels" 
        :key="key"
        :class="{ 'clickable-card': key === 'total_shops' }"
        @click="handleMetricClick(key)"
      >
        <div class="metric-title">{{ val.title }}</div>
        <div class="metric-value">
          {{ formatNumber(overview[key]) }}
          <span class="metric-unit">{{ val.unit }}</span>
        </div>
        <div class="click-hint" v-if="key === 'total_shops'">
          👆 点击查看明细数据
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

    <!-- ================= 弹窗区域 ================= -->

    <!-- 商户列表弹窗 (下钻) -->
    <transition name="fade">
      <div class="modal-overlay" v-if="showShopModal" @click.self="showShopModal = false">
        <div class="modal-content glass-panel">
          <div class="modal-header">
            <h2>🏅 杭州高分商户排行榜</h2>
            <button class="close-btn" @click="showShopModal = false">×</button>
          </div>
          <div class="modal-body">
            <!-- 商户多维筛选工具栏 -->
            <div class="filter-bar">
              <select class="filter-select" v-model="filterDistrict" @change="handleFilterChange">
                <option value="">全部区域</option>
                <option v-for="d in districtOptions" :value="d.id" :key="d.id">{{ d.district_name }}</option>
              </select>
              
              <select class="filter-select" v-model="filterCategory" @change="handleFilterChange">
                <option value="">全部菜系</option>
                <option v-for="c in categoryOptions" :value="c.id" :key="c.id">{{ c.name }}</option>
              </select>
              
              <input class="filter-input" type="text" v-model="filterKeyword" placeholder="输入店名搜索..." @keyup.enter="handleFilterChange" />
              
              <div class="checkbox-group">
                <label><input type="checkbox" v-model="filterFreeParking" @change="handleFilterChange"> 免费停车</label>
                <label><input type="checkbox" v-model="filterReservable" @change="handleFilterChange"> 可订座</label>
                <label><input type="checkbox" v-model="filterBabyChair" @change="handleFilterChange"> 宝宝椅</label>
                <label><input type="checkbox" v-model="filterPrivateRoom" @change="handleFilterChange"> 包厢</label>
              </div>

              <button class="filter-btn" @click="handleFilterChange">🔍 筛选查询</button>
            </div>

            <table class="data-table">
              <thead>
                <tr>
                  <th>排名</th>
                  <th>商户名称</th>
                  <th>所属区域</th>
                  <th>主营菜系</th>
                  <th>人均消费</th>
                  <th>综合评分</th>
                  <th>人气(评价数)</th>
                  <th>服务设施</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(shop, index) in shopList" :key="shop.id">
                  <td><span class="rank-badge" :class="'rank-' + (index + 1)">{{ (shopPage - 1) * shopPageSize + index + 1 }}</span></td>
                  <td class="shop-name">{{ shop.name }}</td>
                  <td>{{ getDistrictName(shop.district_id) }}</td>
                  <td><span class="category-tag">{{ getCategoryName(shop.category_id) }}</span></td>
                  <td class="price-text">¥{{ shop.avg_price }}</td>
                  <td>⭐ {{ shop.rating }}</td>
                  <td>🔥 {{ shop.review_count }}</td>
                  <td class="facilities-cell">
                    <span v-if="shop.has_free_parking" class="f-tag" title="免费停车">🅿️</span>
                    <span v-if="shop.is_reservable" class="f-tag" title="可订座">📅</span>
                    <span v-if="shop.has_baby_chair" class="f-tag" title="宝宝椅">👶</span>
                    <span v-if="shop.has_private_room" class="f-tag" title="包厢">🚪</span>
                  </td>
                </tr>
              </tbody>
            </table>
            
            <div class="pagination">
              <button class="page-btn" :disabled="shopPage <= 1" @click="changeShopPage(-1)">上一页</button>
              <span class="page-info">第 {{ shopPage }} 页 / 共 {{ Math.ceil(shopTotal / shopPageSize) }} 页</span>
              <button class="page-btn" :disabled="shopPage >= Math.ceil(shopTotal / shopPageSize)" @click="changeShopPage(1)">下一页</button>
            </div>
          </div>
        </div>
      </div>
    </transition>



  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import request from '../utils/request'

const router = useRouter()
const username = ref(localStorage.getItem('username') || 'Admin')

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}

// ================= 指标卡数据模型 =================
const metricsLabels = {
  total_shops: { title: '收录商户总数', unit: '家' },
  avg_price: { title: '全市人均消费', unit: '元' },
  avg_rating: { title: '全网平均星级', unit: '星' }
}
const overview = ref({
  total_shops: 0,
  avg_price: 0,
  avg_rating: 0
})

const formatNumber = (num) => {
  if (!num) return 0
  return Number(num).toLocaleString(undefined, { maximumFractionDigits: 1 })
}

// ================= 弹窗交互逻辑 (数据下钻) =================

// 1. 商户列表弹窗
const showShopModal = ref(false)
const shopList = ref([])
const shopPage = ref(1)
const shopPageSize = ref(10)
const shopTotal = ref(0)

// 筛选条件状态
const filterDistrict = ref('')
const filterCategory = ref('')
const filterKeyword = ref('')
const filterFreeParking = ref(false)
const filterReservable = ref(false)
const filterBabyChair = ref(false)
const filterPrivateRoom = ref(false)

const districtOptions = ref([])
const categoryOptions = ref([])

// 加载字典数据 (区域和分类下拉框)
const loadDictionaries = async () => {
  try {
    const [distRes, catRes] = await Promise.all([
      request.get('/districts'),
      request.get('/categories')
    ])
    districtOptions.value = distRes
    categoryOptions.value = catRes
  } catch (e) {
    console.error("加载字典数据失败", e)
  }
}

// 处理卡片点击
const handleMetricClick = (key) => {
  if (key === 'total_shops') {
    showShopModal.value = true
    shopPage.value = 1
    // 每次打开时，如果还没加载字典，就加载一下
    if (districtOptions.value.length === 0) {
      loadDictionaries()
    }
    fetchShops()
  }
}

const handleFilterChange = () => {
  shopPage.value = 1 // 切换筛选条件时，重置到第一页
  fetchShops()
}

// 拉取商户列表数据 (按评分降序 + 动态多维筛选)
const fetchShops = async () => {
  try {
    const res = await request.get('/restaurants', {
      params: { 
        page: shopPage.value, 
        page_size: shopPageSize.value, 
        sort_by: 'rating',
        district_id: filterDistrict.value || undefined,
        category_id: filterCategory.value || undefined,
        keyword: filterKeyword.value || undefined,
        has_free_parking: filterFreeParking.value || undefined,
        is_reservable: filterReservable.value || undefined,
        has_baby_chair: filterBabyChair.value || undefined,
        has_private_room: filterPrivateRoom.value || undefined
      }
    })
    shopList.value = res.list
    shopTotal.value = res.total
  } catch (e) {
    console.error("获取商户列表失败", e)
  }
}
const changeShopPage = (delta) => {
  shopPage.value += delta
  fetchShops()
}

// 简单的字典映射
const districtsMap = {1:'西湖区', 2:'上城区', 3:'拱墅区', 4:'滨江区', 5:'萧山区', 6:'余杭区', 7:'临平区', 8:'钱塘区', 9:'富阳区', 10:'临安区'}
const getDistrictName = (id) => districtsMap[id] || '未知区域'

// ================= ECharts 挂载与渲染逻辑 =================
const pieChartRef = ref(null)
const barChartRef = ref(null)
const scatterChartRef = ref(null)

let pieChart, barChart, scatterChart
const textStyle = { color: '#94a3b8' } 

const loadData = async () => {
  try {
    const [overRes, pieRes, barRes, scatterRes] = await Promise.all([
      request.get('/statistics/overview'),
      request.get('/statistics/category-pie'),
      request.get('/statistics/district-bar'),
      request.get('/statistics/price-rating-scatter')
    ])

    overview.value = overRes
    renderPieChart(pieRes)
    renderBarChart(barRes)
    renderScatterChart(scatterRes)

  } catch (error) {
    console.error("加载大屏数据失败", error)
  }
}

const renderPieChart = (data) => {
  pieChart = echarts.init(pieChartRef.value)
  pieChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b} : {c}家 ({d}%)' },
    series: [
      {
        name: '品类占比',
        type: 'pie',
        radius: [20, 100], 
        center: ['50%', '50%'],
        roseType: 'area',  
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
      min: 2, max: 5, 
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
  position: relative;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border-top: 3px solid var(--brand-primary);
  transition: all var(--transition-fast);
}

.clickable-card {
  cursor: pointer;
  border-top: 3px solid #34d399;
}

.clickable-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 10px 25px rgba(52, 211, 153, 0.2);
  border-color: #10b981;
}

.click-hint {
  position: absolute;
  top: 10px;
  right: 15px;
  font-size: 0.8rem;
  color: #34d399;
  opacity: 0.8;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 0.4; }
  50% { opacity: 1; }
  100% { opacity: 0.4; }
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
  grid-template-columns: 1fr 2.5fr; 
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

/* ================= 弹窗样式 (Modal) ================= */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 800px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  padding: 2rem;
  animation: slideDown 0.3s ease-out forwards;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 1rem;
}

.modal-header h2 {
  font-size: 1.4rem;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 2rem;
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #ef4444;
}

.modal-body {
  overflow-y: auto;
  padding-right: 0.5rem;
}

/* 自定义滚动条 */
.modal-body::-webkit-scrollbar {
  width: 6px;
}
.modal-body::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

/* 数据表格 */
.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1.5rem;
}

.data-table th, .data-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.data-table th {
  color: var(--text-secondary);
  font-weight: 600;
  background: rgba(255, 255, 255, 0.02);
}

.data-table tr:hover {
  background: rgba(255, 255, 255, 0.03);
}

.shop-name {
  font-weight: 600;
  color: #60a5fa;
}

.price-text {
  color: #f59e0b;
  font-family: monospace;
  font-size: 1.1rem;
}

/* 筛选工具栏 */
.filter-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  align-items: center;
  background: rgba(255, 255, 255, 0.02);
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.filter-select, .filter-input {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 0.6rem 1rem;
  border-radius: 6px;
  outline: none;
  font-size: 0.95rem;
  transition: all 0.3s;
}

.filter-select:focus, .filter-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.filter-select option {
  background: #1e293b;
  color: #fff;
}

.filter-btn {
  background: linear-gradient(to right, #3b82f6, #0ea5e9);
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s;
}

.filter-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.checkbox-group {
  display: flex;
  gap: 1rem;
  color: #fff;
  font-size: 0.9rem;
  align-items: center;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  cursor: pointer;
}

.checkbox-group input[type="checkbox"] {
  cursor: pointer;
  accent-color: #3b82f6;
}

.facilities-cell {
  display: flex;
  gap: 0.3rem;
  justify-content: center;
}

.f-tag {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: help;
}


/* 评价瀑布流 */
.review-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.review-item {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  padding: 1.2rem;
  transition: transform 0.2s;
}

.review-item:hover {
  transform: translateX(5px);
  border-color: rgba(52, 211, 153, 0.3);
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.user-name {
  color: #94a3b8;
  font-size: 0.9rem;
}

.review-date {
  color: #64748b;
  font-size: 0.85rem;
}

.review-target {
  font-size: 0.95rem;
  color: #cbd5e1;
  margin-bottom: 0.8rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.target-shop {
  color: #34d399;
  font-weight: 600;
  background: rgba(52, 211, 153, 0.1);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.rating-tag {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: bold;
}

.review-content {
  font-size: 1.05rem;
  line-height: 1.6;
  color: #f8fafc;
  font-style: italic;
}

/* 分页控件 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
  margin-top: 1rem;
}

.page-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:not(:disabled):hover {
  background: var(--brand-primary);
  border-color: var(--brand-primary);
}

.page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-info {
  color: var(--text-secondary);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
