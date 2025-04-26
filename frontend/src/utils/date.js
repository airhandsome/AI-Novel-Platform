// src/utils/date.js

/**
 * 格式化日期
 * @param {Date|string|number} date 日期对象、时间戳或日期字符串
 * @param {string} [format='YYYY-MM-DD HH:mm:ss'] 格式化模式
 * @returns {string} 格式化后的日期字符串
 */
export const formatDate = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
    if (!date) return ''
    
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''
  
    const year = d.getFullYear()
    const month = d.getMonth() + 1
    const day = d.getDate()
    const hours = d.getHours()
    const minutes = d.getMinutes()
    const seconds = d.getSeconds()
  
    const formatMap = {
      'YYYY': year,
      'MM': padZero(month),
      'M': month,
      'DD': padZero(day),
      'D': day,
      'HH': padZero(hours),
      'H': hours,
      'mm': padZero(minutes),
      'm': minutes,
      'ss': padZero(seconds),
      's': seconds
    }
  
    return format.replace(/YYYY|MM|M|DD|D|HH|H|mm|m|ss|s/g, match => formatMap[match])
  }
  
  /**
   * 将数字转为两位数字符串
   * @param {number} num 
   * @returns {string}
   */
  const padZero = (num) => {
    return String(num).padStart(2, '0')
  }
  
  /**
   * 获取相对时间描述
   * @param {Date|string|number} date 日期对象、时间戳或日期字符串
   * @returns {string} 相对时间描述
   */
  export const getRelativeTime = (date) => {
    if (!date) return ''
    
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''
  
    const now = new Date()
    const diff = now.getTime() - d.getTime()
    const diffMinutes = Math.floor(diff / (1000 * 60))
    const diffHours = Math.floor(diff / (1000 * 60 * 60))
    const diffDays = Math.floor(diff / (1000 * 60 * 60 * 24))
  
    if (diffMinutes < 1) {
      return '刚刚'
    } else if (diffMinutes < 60) {
      return `${diffMinutes}分钟前`
    } else if (diffHours < 24) {
      return `${diffHours}小时前`
    } else if (diffDays < 30) {
      return `${diffDays}天前`
    } else {
      return formatDate(date, 'YYYY-MM-DD')
    }
  }
  
  /**
   * 判断是否为同一天
   * @param {Date|string|number} date1 
   * @param {Date|string|number} date2 
   * @returns {boolean}
   */
  export const isSameDay = (date1, date2) => {
    const d1 = new Date(date1)
    const d2 = new Date(date2)
    return d1.getFullYear() === d2.getFullYear() &&
           d1.getMonth() === d2.getMonth() &&
           d1.getDate() === d2.getDate()
  }
  
  /**
   * 获取日期范围描述
   * @param {Date|string|number} startDate 开始日期
   * @param {Date|string|number} endDate 结束日期
   * @returns {string} 日期范围描述
   */
  export const getDateRange = (startDate, endDate) => {
    if (!startDate || !endDate) return ''
    
    const start = new Date(startDate)
    const end = new Date(endDate)
    
    if (isNaN(start.getTime()) || isNaN(end.getTime())) return ''
    
    if (isSameDay(start, end)) {
      return formatDate(start, 'YYYY-MM-DD')
    }
    
    return `${formatDate(start, 'YYYY-MM-DD')} 至 ${formatDate(end, 'YYYY-MM-DD')}`
  }