<template>
  <div class="page-footer">
    <span>名称: {{ copyright.name }}</span>
    <span class="split">|</span>
    <span>版本: {{ copyright.version }}</span>
    <span class="split">|</span>
    <span>Go版本: {{ copyright.goVersion }}</span>
    <span class="split">|</span>
    <span>平台: {{ copyright.platform }}</span>
    <span class="split">|</span>
    <span>PID: {{ copyright.pid }}</span>
    <span class="split">|</span>
    <span>编译时间: {{ copyright.buildTime }}</span>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { getCopyright } from '@/api/dns'

const copyright = ref({
  name: '本地DNS代理',
  version: '1.0.0',
  platform: 'windows/amd64',
  goVersion: '1.0.0',
  pid: '0',
  buildTime: '2024-01-01 00:00:00'
})

getCopyright().then(res => {
  copyright.value = res.data
})
</script>

<style lang="scss" scoped>
.page-footer {
  text-align: center;
  padding: 8px;
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: center;
  background-color: var(--page-layout-background-color);

  .split {
    user-select: none;
    color: var(--vxe-ui-status-info-color);
  }
}
</style>
