import { VxeUI } from 'vxe-pc-ui'

// 全局参数
VxeUI.setConfig({
  version: 0,
  zIndex: 999,

  table: {
    border: true,
    showOverflow: true,
    autoResize: true,
    columnConfig: {
      resizable: true
    },
    editConfig: {
      trigger: 'click'
    },
    sortConfig: {
      trigger: 'cell'
    },
    scrollY: {
      enabled: true,
      gt: 20
    }
  },
  grid: {
    toolbarConfig: {
      custom: true
    },
    proxyConfig: {
      showResponseMsg: false,
      showActiveMsg: true,
      response: {
        total: 'page.total',
        result: 'data',
        list: 'data'
      },
      ajax: {
        deleteSuccess () {
          VxeUI.modal.message({
            content: '删除成功',
            status: 'success'
          })
        },
        saveSuccess ({ response }) {
          const { data } = response
          VxeUI.modal.message({
            content: `新增 ${data.insertCount} 条，删除 ${data.deleteCount} 条，修改 ${data.updateCount} 条`,
            status: 'success'
          })
        }
      }
    }
  },
  pager: {
    layouts: ['Home', 'PrevJump', 'PrevPage', 'Jump', 'PageCount', 'NextPage', 'NextJump', 'End', 'Sizes', 'Total']
  }
})
