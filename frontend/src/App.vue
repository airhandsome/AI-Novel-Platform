<template>
  <div class="app">
    <NavBar />
    <div class="main-content">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import NavBar from './components/layout/NavBar.vue'
</script>

<style>
.app {
  min-height: 100vh;
}

.main-content {
  padding-top: 60px; /* 与导航栏高度相同 */
  min-height: calc(100vh - 60px);
}

/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background: #f5f7fa;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style> 