import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import { getUpcomingDoses } from "@/api/doses";
import { cacheGet, cacheSet } from "@/utils/cache";

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

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>Upcoming</Text>

      <FlatList
        data={doses}
        keyExtractor={(d) => d.id}
        renderItem={({ item }) => (
          <DoseCard
            dose={item}
            onTake={() => {}}
            onSkip={() => {}}
          />
        )}
      />
    </View>
  );
}
