<template>
  <div class="login-container">
    <div class="glass-panel login-box">
      <h2 class="title">{{ isLogin ? '系统登录' : '创建账号' }}</h2>
      <p class="subtitle">杭州餐饮大数据分析平台</p>

      <form @submit.prevent="handleSubmit" class="form-container">
        <div class="input-group">
          <input 
            v-model="form.username" 
            type="text" 
            class="input-base" 
            placeholder="请输入账号 (至少3位)" 
            required 
            minlength="3"
          />
        </div>
        
        <div class="input-group">
          <input 
            v-model="form.password" 
            type="password" 
            class="input-base" 
            placeholder="请输入密码 (至少6位)" 
            required 
            minlength="6"
          />
        </div>

        <button type="submit" class="btn submit-btn" :disabled="loading">
          {{ loading ? '处理中...' : (isLogin ? '登 录' : '注 册') }}
        </button>

        <p class="toggle-text" @click="toggleMode">
          {{ isLogin ? '没有账号？立即注册' : '已有账号？返回登录' }}
        </p>
      </form>
      
      <!-- 错误提示组件，带有淡入淡出动画 -->
      <transition name="fade">
        <div v-if="errorMsg" class="error-msg">
          {{ errorMsg }}
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'

const router = useRouter()

const isLogin = ref(true)
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: ''
})

// 切换登录/注册模式
const toggleMode = () => {
  isLogin.value = !isLogin.value
  errorMsg.value = ''
  form.username = ''
  form.password = ''
}

// 提交表单
const handleSubmit = async () => {
  errorMsg.value = ''
  loading.value = true
  
  try {
    if (isLogin.value) {
      // 发送登录请求
      const res = await request.post('/user/login', form)
      
      // 保存 Token 和用户名到浏览器缓存
      localStorage.setItem('token', res.token)
      localStorage.setItem('username', res.username)
      
      // 成功后跳转到数据大屏
      router.push('/dashboard')
    } else {
      // 发送注册请求
      await request.post('/user/register', form)
      alert('注册成功，请进行登录！')
      
      // 注册成功后自动切回登录页面并清空密码
      isLogin.value = true
      form.password = ''
    }
  } catch (err) {
    // 捕获 axios 拦截器里抛出的报错信息
    errorMsg.value = err.message || '网络连接异常，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
}

.login-box {
  width: 100%;
  max-width: 420px;
  padding: 3.5rem 3rem;
  text-align: center;
  /* 高级微浮现入场动画 */
  animation: floatUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  transform: translateY(30px);
  opacity: 0;
}

@keyframes floatUp {
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  /* 炫酷的文字渐变效果 */
  background: linear-gradient(to right, #ffffff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 0.95rem;
  margin-bottom: 2.5rem;
  letter-spacing: 1px;
}

.form-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-group {
  position: relative;
}

.submit-btn {
  margin-top: 1rem;
  padding: 0.85rem;
  font-size: 1.1rem;
  letter-spacing: 4px;
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.toggle-text {
  margin-top: 1rem;
  color: var(--text-muted);
  font-size: 0.9rem;
  cursor: pointer;
  transition: var(--transition-fast);
}

.toggle-text:hover {
  color: var(--brand-accent);
}

.error-msg {
  margin-top: 1.5rem;
  padding: 0.75rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  color: #fca5a5;
  border-radius: 8px;
  font-size: 0.9rem;
}

/* 错误提示框的渐变动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
