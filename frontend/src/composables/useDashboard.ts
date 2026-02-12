import { ref, computed } from 'vue'

export interface DashboardStats {
  totalUsers: number
  activeUsers: number
  newRegistrations: number
  activeSessions: number
}

export interface ActivityItem {
  id: string
  user: string
  action: string
  time: string
  type: 'user' | 'create' | 'delete' | 'update' | 'security'
  changes?: string
}

export interface ChartData {
  label: string
  value: number
}

export function useDashboard() {
  const loading = ref(true)
  const stats = ref<DashboardStats>({
    totalUsers: 0,
    activeUsers: 0,
    newRegistrations: 0,
    activeSessions: 0
  })

  const activities = ref<ActivityItem[]>([])
  const chartData = ref<ChartData[]>([])

  // Simulated data - replace with actual API calls
  const loadDashboardData = async () => {
    loading.value = true

    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 800))

    // Mock data
    stats.value = {
      totalUsers: 2847,
      activeUsers: 1923,
      newRegistrations: 47,
      activeSessions: 312
    }

    activities.value = [
      {
        id: '1',
        user: 'John Doe',
        action: 'created a new account',
        time: '2 minutes ago',
        type: 'create'
      },
      {
        id: '2',
        user: 'Sarah Wilson',
        action: 'updated profile settings',
        time: '5 minutes ago',
        type: 'update',
        changes: '3 fields'
      },
      {
        id: '3',
        user: 'Mike Johnson',
        action: 'logged in from new device',
        time: '12 minutes ago',
        type: 'security'
      },
      {
        id: '4',
        user: 'Admin',
        action: 'deleted user account',
        time: '25 minutes ago',
        type: 'delete'
      },
      {
        id: '5',
        user: 'Emily Brown',
        action: 'updated email address',
        time: '1 hour ago',
        type: 'update'
      },
      {
        id: '6',
        user: 'David Lee',
        action: 'changed password',
        time: '2 hours ago',
        type: 'security'
      }
    ]

    // Get last 7 days for chart
    const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    const today = new Date().getDay()
    const last7Days = []

    for (let i = 6; i >= 0; i--) {
      const dayIndex = (today - i + 7) % 7
      last7Days.push({
        label: days[dayIndex],
        value: Math.floor(Math.random() * 50) + 20
      })
    }

    chartData.value = last7Days

    loading.value = false
  }

  const refreshData = () => {
    loadDashboardData()
  }

  // Derived values for stat cards
  const totalUsersProgress = computed(() => {
    return Math.min((stats.value.totalUsers / 5000) * 100, 100)
  })

  const activeUsersProgress = computed(() => {
    return Math.min((stats.value.activeUsers / stats.value.totalUsers) * 100, 100)
  })

  const newRegistrationsProgress = computed(() => {
    return Math.min((stats.value.newRegistrations / 100) * 100, 100)
  })

  const activeSessionsProgress = computed(() => {
    return Math.min((stats.value.activeSessions / 500) * 100, 100)
  })

  return {
    loading,
    stats,
    activities,
    chartData,
    totalUsersProgress,
    activeUsersProgress,
    newRegistrationsProgress,
    activeSessionsProgress,
    loadDashboardData,
    refreshData
  }
}
