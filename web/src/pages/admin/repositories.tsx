import { useAuth } from "@/auth/AuthContext";
import SwirlingEffectSpinner from "@/components/customized/spinner/spinner-06";
import { useGetCronInfoQuery } from "@/generated/graphql"
import { ResponsiveHeatMap } from '@nivo/heatmap'
import { LogIn } from "lucide-react";
import { useEffect, useState } from "react"
import { RequestLogin } from '@/components/request-login'
import { toast } from "sonner";

type HeatmapEntry = {
  id: string;
  data: {
    x: string;
    y: number;
  }[];
}

export const Repositories = () => {
  const { jwt } = useAuth()
  const { data, loading, error } = useGetCronInfoQuery()
  const [heatmapData, setHeatmapData] = useState<HeatmapEntry[]>([])

  const cronInfo = data?.getCronInfo ?? []
  const allDays = [
    "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
  ];
  const allHours = Array.from({ length: 24 }, (_, i) => i.toString());


  useEffect(() => {
    if (loading || !data?.getCronInfo || heatmapData.length > 0) {
      return
    }

    const inputMap = new Map(cronInfo.map(cron => [`${cron.day}-${cron.hour}`, cron.total]));
    const result = allDays.map(day => ({
      id: day,
      data: allHours.map(hour => ({
        x: hour,
        y: inputMap.get(`${day}-${hour}`) ?? 0
      }))
    }));

    setHeatmapData(result)
  }, [data])


  useEffect(() => {
    if (!error) {
      return
    }

    toast.error("Could not fetch dependencies", {
      description: error.message,
    })

  }, [error])

  return (
    <div>
      {!jwt && (
        <RequestLogin />
      )}
      {loading && jwt && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 z-50">
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
              maxValue: 10
            }}
          />
        )}
      </div>
    </div >
  )
}
