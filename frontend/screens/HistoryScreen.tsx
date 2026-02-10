import { useEffect, useState } from "react";
import { View, Text, ScrollView } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import { API } from "@/api/doses";

export default function HistoryScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);

  type GroupedDoses = Record<string, { 
    label: string; 
    times: Record<string, (Dose & { displayTime: string })[]> 
  }>;

  useEffect(() => {
    const fetchHistoryDoses = async () => {
      const res = await API.get<Dose[]>("/doses/history");
      setDoses(res.data);
    };
    fetchHistoryDoses();
  }, []);

  // Group doses by date -> time
  const groupByDateAndTime = (doses: Dose[]): GroupedDoses => {
    const map: GroupedDoses = {};

    doses.forEach((dose) => {
      const sourceDate = dose.taken_at ? new Date(dose.taken_at) : new Date(dose.scheduled_at);
      if (isNaN(sourceDate.getTime())) return;

      const dateKey = sourceDate.toISOString().slice(0, 10); // YYYY-MM-DD
      const timeKey = sourceDate.toISOString().slice(11, 16); // HH:mm

      const dateLabel = sourceDate.toLocaleDateString([], {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });

      const timeLabel = sourceDate.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        hour12: true,
      });

      // Initialize date entry
      if (!map[dateKey]) map[dateKey] = { label: dateLabel, times: {} };

      // Initialize time entry
      if (!map[dateKey].times[timeKey]) map[dateKey].times[timeKey] = [];

      // Push dose with displayTime
      map[dateKey].times[timeKey].push({ ...dose, displayTime: timeLabel });
    });

    return map;
  };

  const grouped = groupByDateAndTime(doses);

  // Convert grouped object into sorted array for rendering
  const groupByDatesData = Object.entries(grouped)
    // Sort dates descending (newest first)
    .sort(([a], [b]) => b.localeCompare(a))
    .map(([_, { label, times }]) => [
      label,
      Object.entries(times)
        // Sort times ascending
        .sort(([t1], [t2]) => t1.localeCompare(t2))
        // Keep both timeLabel and doses array
        .map(([timeLabel, dosesAtTime]) => [timeLabel, dosesAtTime] as const),
    ] as const);

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>History</Text>
      <ScrollView>
        {groupByDatesData.map(([dateLabel, times]) => (
          <View key={dateLabel} style={{ marginBottom: 24 }}>
            <Text style={{ fontSize: 18, fontWeight: "600", marginBottom: 10 }}>
              {dateLabel}
            </Text>

            {times.map(([timeLabel, dosesAtTime]) => (
              <View key={timeLabel} style={{ marginBottom: 16 }}>
                <Text style={{ fontSize: 14, fontWeight: "500", marginBottom: 6 }}>
                  {dosesAtTime[0]?.displayTime || timeLabel}
                </Text>

                {dosesAtTime.map((dose) => (
                  <DoseCard
                    key={dose.id}
                    dose={dose}
                    onSkip={() => {}}
                    onTake={() => {}}
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
