// Export utility composable
export const useExport = () => {
  // Convert data to CSV string
  const toCSV = (data: any[], columns: { key: string; label: string }[]) => {
    if (!data.length) return ''
    
    // Header row
    const headers = columns.map(col => `"${col.label}"`).join(',')
    
    // Data rows
    const rows = data.map(item => {
      return columns.map(col => {
        let value = getNestedValue(item, col.key)
        
        // Handle different types
        if (value === null || value === undefined) {
          value = ''
        } else if (typeof value === 'object') {
          value = JSON.stringify(value)
        } else if (typeof value === 'boolean') {
          value = value ? 'Ya' : 'Tidak'
        }
        
        // Escape quotes and wrap in quotes
        return `"${String(value).replace(/"/g, '""')}"`
      }).join(',')
    })
    
    return [headers, ...rows].join('\n')
  }
  
  // Get nested object value (e.g., "user.name")
  const getNestedValue = (obj: any, path: string) => {
    return path.split('.').reduce((acc, part) => acc && acc[part], obj)
  }
  
  // Download CSV file
  const downloadCSV = (csvContent: string, filename: string) => {
    // Add BOM for Excel compatibility with UTF-8
    const BOM = '\uFEFF'
    const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    
    if (navigator.msSaveBlob) {
      // IE 10+
      navigator.msSaveBlob(blob, filename)
    } else {
      const url = URL.createObjectURL(blob)
      link.href = url
      link.download = filename
      link.style.display = 'none'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)
    }
  }
  
  // Export users
  const exportUsers = (users: any[]) => {
    const columns = [
      { key: 'id', label: 'ID' },
      { key: 'full_name', label: 'Nama Lengkap' },
      { key: 'email', label: 'Email' },
      { key: 'role', label: 'Role' },
      { key: 'is_active', label: 'Status Aktif' },
      { key: 'enrollment_count', label: 'Jumlah Kursus' },
      { key: 'created_at', label: 'Tanggal Daftar' }
    ]
    
    const csv = toCSV(users, columns)
    const date = new Date().toISOString().split('T')[0]
    downloadCSV(csv, `users_${date}.csv`)
  }
  
  // Export courses
  const exportCourses = (courses: any[]) => {
    const columns = [
      { key: 'id', label: 'ID' },
      { key: 'title', label: 'Judul' },
      { key: 'category.name', label: 'Kategori' },
      { key: 'instructor.full_name', label: 'Instruktur' },
      { key: 'price', label: 'Harga' },
      { key: 'is_published', label: 'Published' },
      { key: 'created_at', label: 'Tanggal Dibuat' }
    ]
    
    const csv = toCSV(courses, columns)
    const date = new Date().toISOString().split('T')[0]
    downloadCSV(csv, `courses_${date}.csv`)
  }
  
  // Export transactions
  const exportTransactions = (transactions: any[]) => {
    const columns = [
      { key: 'id', label: 'ID Transaksi' },
      { key: 'user.full_name', label: 'Nama User' },
      { key: 'user.email', label: 'Email User' },
      { key: 'course.title', label: 'Kursus' },
      { key: 'amount', label: 'Jumlah' },
      { key: 'currency', label: 'Mata Uang' },
      { key: 'status', label: 'Status' },
      { key: 'payment_gateway', label: 'Gateway' },
      { key: 'created_at', label: 'Tanggal' }
    ]
    
    const csv = toCSV(transactions, columns)
    const date = new Date().toISOString().split('T')[0]
    downloadCSV(csv, `transactions_${date}.csv`)
  }
  
  // Export instructors
  const exportInstructors = (instructors: any[]) => {
    const columns = [
      { key: 'id', label: 'ID' },
      { key: 'full_name', label: 'Nama Lengkap' },
      { key: 'email', label: 'Email' },
      { key: 'specialty', label: 'Spesialisasi' },
      { key: 'course_count', label: 'Jumlah Kursus' },
      { key: 'is_active', label: 'Status Aktif' }
    ]
    
    const csv = toCSV(instructors, columns)
    const date = new Date().toISOString().split('T')[0]
    downloadCSV(csv, `instructors_${date}.csv`)
  }
  
  // Generic export
  const exportData = (data: any[], columns: { key: string; label: string }[], filename: string) => {
    const csv = toCSV(data, columns)
    downloadCSV(csv, filename)
  }
  
  return {
    toCSV,
    downloadCSV,
    exportUsers,
    exportCourses,
    exportTransactions,
    exportInstructors,
    exportData
  }
}

