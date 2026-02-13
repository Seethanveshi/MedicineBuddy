import { useEffect, useState } from "react";
import { View, Text, ScrollView } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import { API } from "@/api/doses";
import { format } from "date-fns";

export default function HistoryScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);

  useEffect(() => {
    const fetchHistoryDoses = async () => {
      const res = await API.get<Dose[]>("/doses/history");
      setDoses(res.data);
    };
    fetchHistoryDoses();
  }, []);

  const groupByDate = (doses: Dose[]) => {
    const map: Record<string, Dose[]> = {};

    doses.forEach((dose) => {
      const date = new Date(dose.scheduled_at).toISOString().split("T")[0];

      if (!map[date]) map[date] = [];
      map[date].push(dose);
    });

    return Object.entries(map).sort(
      ([a], [b]) => new Date(b).getTime() - new Date(a).getTime()
    );
  };

  const groupByTime = (doses: Dose[]) => {
    const map: Record<string, Dose[]> = {};

    doses.forEach((dose) => {
      const time = new Date(dose.scheduled_at).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });

      if (!map[time]) map[time] = [];
      map[time].push(dose);
    });

    return Object.entries(map).sort(
      ([a], [b]) =>
        Date.parse(`1970-01-01 ${a}`) - Date.parse(`1970-01-01 ${b}`)
    );

  };

  const dateGroups = groupByDate(doses);

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>History</Text>
      <ScrollView>
          {dateGroups.map(([date, dateDoses]) => (
              <View key={date} style={{ marginBottom: 24 }}>
                <Text style={{
                  fontSize: 18,
                  fontWeight: "700",
                  marginBottom: 12,
                }}>
                  {format(new Date(date), "MMM dd")}
                </Text>
    
                {/* TIME GROUPS */}
                {groupByTime(dateDoses).map(([time, timeDoses]) => (
                  <View key={time} style={{ marginBottom: 16 }}>
                    <Text style={{
                      fontSize: 16,
                      fontWeight: "600",
                      marginBottom: 8,
                    }}>
                      {time}
                    </Text>
    
                    {timeDoses.map((dose) => (
                      <DoseCard
                        key={dose.id}
                        dose={dose}
                        onTake={async () => {}}
                        onSkip={async () => {}}
                      />
                    ))}
                  </View>
                ))}
              </View>
            ))}
      </ScrollView>
    </View>
  );
}
