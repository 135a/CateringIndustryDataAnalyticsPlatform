<script setup>
// App 组件只需提供一个 router-view 出口即可
</script>

<template>
  <router-view v-slot="{ Component }">
    <!-- 引入页面切换时的全局淡入淡出动画 -->
    <transition name="fade" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<style scoped>
/* Vue 官方的 Transition 动画控制类 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
