import { useAuth } from '@/auth/AuthContext'
import SwirlingEffectSpinner from '@/components/customized/spinner/spinner-06'
import { useGetCronInfoQuery } from '@/generated/graphql'
import { ResponsiveHeatMap } from '@nivo/heatmap'
import { useEffect, useState } from 'react'
import { RequestLogin } from '@/components/RequestLogin'
import { toast } from 'sonner'
import { useTheme } from '@/components/ThemeProvider'

type HeatmapEntry = {
  id: string
  data: {
    x: string
    y: number
  }[]
}

export const Repositories = () => {
  const { jwt } = useAuth()
  const { data, loading, error } = useGetCronInfoQuery()
  const [heatmapData, setHeatmapData] = useState<HeatmapEntry[]>([])
  const { theme } = useTheme()

  useEffect(() => {
    if (loading || !data?.getCronInfo || heatmapData.length > 0) {
      return
    }

    const cronInfo = data?.getCronInfo ?? []
    const allDays = [
      'Sunday',
      'Monday',
      'Tuesday',
      'Wednesday',
      'Thursday',
      'Friday',
      'Saturday',
    ]
    const allHours = Array.from({ length: 24 }, (_, i) => i.toString())

    const inputMap = new Map(
      cronInfo.map((cron) => [`${cron.day}-${cron.hour}`, cron.total]),
    )
    const result = allDays.map((day) => ({
      id: day,
      data: allHours.map((hour) => ({
        x: (+hour + 1).toString(),
        y: inputMap.get(`${day}-${hour}`) ?? 0,
      })),
    }))

    setHeatmapData(result)
  }, [data, heatmapData.length, loading])

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error('Could not fetch dependencies', {
      description: error.message,
    })
  }, [error])

  const labelColor = theme === 'dark' ? '#E5E7EB' : '#1F2937' // tailwind slate-200 / slate-800

  return (
    <div>
      {!jwt && <RequestLogin />}
      {loading && jwt && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 dark:bg-slate-900/80 z-50">
          <SwirlingEffectSpinner />
        </div>
      )}
      <div style={{ height: 500 }}>
        {jwt && !loading && heatmapData.length > 0 && (
          <ResponsiveHeatMap
            data={heatmapData}
            margin={{ top: 60, right: 0, bottom: 60, left: 100 }}
            emptyColor="#555555"
            colors={{
              type: 'diverging',
              scheme: 'greens',
              divergeAt: 0.5,
              minValue: 0,
              maxValue: 10,
            }}
            axisTop={{
              tickSize: 5,
              tickPadding: 5,
              tickRotation: -45,
              legend: '',
              legendOffset: 46,
              tickValues: 'every 1',
            }}
            axisLeft={{
              tickSize: 5,
              tickPadding: 5,
              tickRotation: 0,
              legend: '',
              legendOffset: -60,
              tickValues: 'every 1',
            }}
            theme={{
              axis: {
                ticks: {
                  text: {
                    fill: labelColor,
                  },
                },
              },
              labels: {
                text: {
                  fill: labelColor,
                },
              },
            }}
          />
        )}
      </div>
    </div>
  )
}
