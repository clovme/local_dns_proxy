export interface UserMenuVO {
  title: string
  code: string
  name: string
  parentCode: string
  routeName: string
  icon: string
  type: string
  routerLink: any
  menuType: 'route' | 'form'
}

export interface UserRouteConfigVO {
  name: string
  code: string
  parentCode: string
  routeName: string
  type: 'menu' | 'action'
  icon: string
  sort: number
}
