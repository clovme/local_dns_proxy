<template>
  <PageView>
    <vxe-grid ref="gridRef" v-bind="gridOptions">
      <template #domain="{ row }">
        <vxe-link v-if="row.status === 1" :href="`${row.protocol}://${row.domain}:${row.port}`" target="_blank">{{ row.domain }}</vxe-link>
        <vxe-text v-else>{{ row.domain }}</vxe-text>
      </template>
      <template #description="{ row }">
        本地域名
        <vxe-link v-if="row.status === 1" :href="`${row.protocol}://${row.domain}:${row.port}`" target="_blank">{{ row.domain }}</vxe-link>
        <vxe-text v-else>{{ row.domain }}</vxe-text>
        映射到 {{ row.ip }}
      </template>
      <template #action="{ row }">
        <vxe-button mode="text" status="error" icon="vxe-icon-delete" @click="removeRow(row)">删除</vxe-button>
      </template>
    </vxe-grid>
  </PageView>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { VxeGridInstance, VxeGridProps } from 'vxe-table'
import { VxeUI } from 'vxe-pc-ui'
import { DnsVO, getDnsListPage, postDnsSaveBatch, deleteDnsDelete } from '@/api/dns'

const gridRef = ref<VxeGridInstance<DnsVO>>()

const statusCellRender = {
  name: 'VxeSwitch',
  props: {
    size: 'mini',
    openLabel: '启用',
    closeLabel: '停用',
    openValue: 1,
    closeValue: 0
  }
}

const protocolCellRender = {
  name: 'VxeSwitch',
  props: {
    size: 'mini',
    openLabel: 'HTTPS',
    closeLabel: 'HTTP',
    openValue: 'https',
    closeValue: 'http',
    class: 'protocol-switch'
  }
}

const gridOptions = reactive<VxeGridProps<DnsVO>>({
  id: 'DnsList',
  height: '100%',
  stripe: true,
  keepSource: true,
  showOverflow: true,
  rowConfig: {
    isHover: true
  },
  sortConfig: {
    remote: true,
    multiple: true
  },
  editConfig: {
    mode: 'row',
    showStatus: true,
    trigger: 'dblclick'
  },
  editRules: {
    name: [
      { required: true, message: '请输入名称' }
    ]
  },
  customConfig: {
    storage: true
  },
  toolbarConfig: {
    refresh: true,
    zoom: true,
    buttons: [
      { name: '新增', code: 'insert_edit', status: 'primary', icon: 'vxe-icon-add' },
      { name: '标记/删除', code: 'mark_cancel', status: 'error', icon: 'vxe-icon-delete' },
      { name: '保存', code: 'save', status: 'success', icon: 'vxe-icon-save' }
    ]
  },
  pagerConfig: {
    enabled: false
  },
  columns: [
    { type: 'checkbox', width: 60 },
    { type: 'seq', width: 60 },
    { field: 'protocol', title: 'HTTP协议', sortable: true, width: 110, cellRender: protocolCellRender },
    { field: 'domain', title: '本地映射域名', sortable: true, width: 300, slots: { default: 'domain' }, editRender: { name: 'VxeInput' } },
    { field: 'ip', title: '映射IP', sortable: true, width: 150, editRender: { name: 'VxeInput' } },
    { field: 'port', title: '映射端口', sortable: true, width: 120, editRender: { name: 'VxeInput' } },
    { field: 'status', title: '启用状态', sortable: true, width: 110, align: 'center', cellRender: statusCellRender },
    { field: 'description', title: '描述', minWidth: 200, slots: { default: 'description' } },
    { field: 'updatedAt', title: '更新时间', width: 160, formatter: 'FormatDateTime', sortable: true },
    { field: 'createdAt', title: '创建时间', width: 160, formatter: 'FormatDateTime', sortable: true },
    { field: 'action', title: '操作', width: 160, slots: { default: 'action' } }
  ],
  proxyConfig: {
    form: false,
    sort: true,
    ajax: {
      query ({ sorts }) {
        const params = {
          orderBy: sorts.map(item => `${item.field}|${item.order}`).join(',')
        }
        return getDnsListPage(params)
      },
      save ({ body }) {
        return postDnsSaveBatch({
          ...body
        })
      }
    }
  }
})

const removeRow = async (row: DnsVO) => {
  const $grid = gridRef.value
  if ($grid && $grid.isInsertByRow(row)) {
    $grid.remove(row)
    return
  }
  const type = await VxeUI.modal.confirm({
    content: `请确认是否删除 “ ${row.domain} ”？`
  })
  if (type === 'confirm') {
    deleteDnsDelete({ id: row.id }).then((res) => {
      if ($grid) {
        $grid.commitProxy('query')
      }
      VxeUI.modal.message({
        content: res.message,
        status: 'success'
      })
    })
  }
}
</script>

<style lang="scss" scoped>
:deep(.protocol-switch.is--on) {
  .vxe-switch--button {
    background-color: var(--vxe-ui-status-success-lighten-color);
  }
}
:deep(.protocol-switch.is--off) {
  .vxe-switch--button {
    background-color: var(--vxe-ui-switch-open-background-color);
  }
}
</style>
