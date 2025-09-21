<template>
  <div class="page-header">
    <div class="header-left">
      <div class="aside-logo">
        <img class="logo-img" draggable="false" src="@/assets/logo.png" />
        <vxe-text class="logo-title">本地DNS代理管理</vxe-text>
      </div>
    </div>
    <div class="header-right">
      <span class="right-item">
        <span class="right-item-title">代理网卡：</span>
        <vxe-pulldown :options="nif" trigger="click" class="right-item-comp" show-popup-shadow transfer  @option-click="networkInterfacesClickEvent">
          <template #default>
            <vxe-text>{{ iface }}</vxe-text>
            <vxe-icon name="caret-down"></vxe-icon>
          </template>
        </vxe-pulldown>
      </span>

      <span class="right-item">
        <span class="right-item-title">服务状态：</span>
        <vxe-switch class="right-item-comp switch-service" v-model="switchService" size="mini" open-value="running" open-label="运行中" close-value="stop" close-label="未启动"></vxe-switch>
      </span>

      <span class="right-item">
        <vxe-switch class="right-item-comp" v-model="currTheme" size="mini" open-value="light" open-label="白天" close-value="dark" close-label="夜间"></vxe-switch>
      </span>

      <span class="right-item">
        <vxe-color-picker class="switch-primary-color" v-model="currPrimaryColor" :colors="colorList" size="mini"></vxe-color-picker>
      </span>

      <span class="right-item">
        <vxe-radio-group class="switch-size" v-model="currCompSize" :options="sizeOptions" type="button" size="mini"></vxe-radio-group>
      </span>

      <span class="right-item">
        <vxe-pulldown :options="langPullList" trigger="click" class="right-item-comp" show-popup-shadow transfer  @option-click="langOptionClickEvent">
          <vxe-button mode="text" icon="vxe-icon-language-switch" :content="langLabel"></vxe-button>
        </vxe-pulldown>
      </span>

      <span class="right-item">
        <div class="user-avatar">
          <vxe-text>DNS</vxe-text>
          <img class="user-picture" draggable="false" src="@/assets/default-picture.gif">
        </div>
      </span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { VxeGlobalI18nLocale, VxePulldownEvents } from 'vxe-pc-ui'
import { useAppStore } from '@/store/app'
import { getNetworkInterfaces } from '@/api/dns'

const iface = ref('')
const nif = ref<{ label: string, value: string }[]>([])

const appStore = useAppStore()

getNetworkInterfaces().then(res => {
  // @ts-ignore
  window.isFirstLoading = 'first'
  const items: { label: string, value: string }[] = []
  res.data.ifaces.forEach(item => {
    items.push({ label: item.name, value: item.name })
  })
  nif.value = items
  iface.value = res.data.iface
  switchService.value = res.data.running
})

const langPullList = ref([
  { label: '中文', value: 'zh-CN' },
  { label: '英文', value: 'en-US' }
])

const switchService = computed({
  get () {
    return appStore.serviceStatus
  },
  set (status) {
    appStore.setServiceStatus(status, iface.value)
  }
})

const langLabel = computed(() => {
  const item = langPullList.value.find(item => item.value === appStore.language)
  return item ? item.label : appStore.language
})

const currTheme = computed({
  get () {
    return appStore.theme
  },
  set (name) {
    appStore.setTheme(name)
  }
})

const currPrimaryColor = computed({
  get () {
    return appStore.primaryColor
  },
  set (color) {
    appStore.setPrimaryColor(color || '')
  }
})

const currCompSize = computed({
  get () {
    return appStore.componentsSize
  },
  set (size) {
    appStore.setComponentsSize(size)
  }
})

const colorList = ref([
  '#409eff', '#29D2F8', '#31FC49', '#3FF2B3', '#B52DFE', '#FC3243', '#FA3077', '#D1FC44', '#FEE529', '#FA9A2C'
])

const sizeOptions = ref([
  { label: '默认', value: '' },
  { label: '中', value: 'medium' },
  { label: '小', value: 'small' },
  { label: '迷你', value: 'mini' }
])

const langOptionClickEvent: VxePulldownEvents.OptionClick = ({ option }) => {
  appStore.setLanguage(option.value as VxeGlobalI18nLocale)
}

const networkInterfacesClickEvent: VxePulldownEvents.OptionClick = ({ option }) => {
  iface.value = option.value as string
}
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  height: 50px;
  padding: 0 16px;
  border-bottom: 1px solid var(--page-layout-border-color);

  .header-left {
    flex-grow: 1;

    .aside-logo {
      display: flex;
      flex-direction: row;
      align-items: center;
      flex-shrink: 0;
      padding: 8px 16px;
      user-select: none;

      .logo-img {
        display: block;
        width: 30px;
        height: 30px;
      }
      .logo-title {
        padding-left: 8px;
        font-weight: 700;
        font-size: 18px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }

  .header-right {
    display: flex;
    flex-direction: row;
    flex-shrink: 0;
    align-items: center;
  }

  .right-item {
    cursor: pointer;
    margin-left: 24px;
  }
  .right-item-title {
    vertical-align: middle;
  }

  .right-item-comp {
    vertical-align: middle;
  }

  .user-avatar {
    display: inline-flex;
    flex-direction: row;
    align-items: center;
    cursor: pointer;
  }

  .user-picture {
    width: 35px;
    height: 35px;
    margin: 0 2px;
  }

  .collapseBtn {
    font-size: 18px;
  }

  .switch-service.vxe-switch.is--on {
    :deep(.vxe-switch--button) {
      background-color: var(--vxe-ui-status-error-color);
    }
  }
}
</style>
