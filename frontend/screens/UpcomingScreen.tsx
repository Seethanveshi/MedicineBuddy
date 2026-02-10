import { useEffect, useState } from "react";
import { View, FlatList, Text, ScrollView } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import { getUpcomingDoses } from "@/api/doses";
import { cacheGet, cacheSet } from "@/utils/cache";
import { format } from "date-fns";
import { groupByDate, groupByTime } from "@/utils/groupBy";

export default function UpcomingScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);
  const UPCOMING_CACHE_KEY = "upcoming_doses";

  useEffect(() => {
    const fetchDoses = async () => {
      // 1️⃣ Load cached upcoming doses
      const cached = await cacheGet<Dose[]>(UPCOMING_CACHE_KEY);
      if (cached) setDoses(cached);

      // 2️⃣ Fetch fresh upcoming doses
      try {
        const data = await getUpcomingDoses();
        setDoses(data);
        await cacheSet(UPCOMING_CACHE_KEY, data);
      } catch (e) {
        console.warn("Failed to fetch upcoming doses, using cache", e);
      }
    };

    fetchDoses();
  }, []);

  const dateGroups = groupByDate(doses);


  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>Upcoming</Text>
      <ScrollView >
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
