  export const formatDate = (dateString: string | null) => {
    if (!dateString) {
      return '-'
    }

    const date = new Date(dateString)
    if (!date) {
      return '-'
    }

    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    })
  }
