import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { getTodayDoses, takeDose, skipDose } from "../api/doses";
import DoseCard from "../components/DoseCard";
import { Dose } from "../types/dose";

export default function TodayScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);

  const load = async () => {
    const data = await getTodayDoses();
    setDoses(data);
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>Today</Text>

      <FlatList
        data={doses}
        keyExtractor={(d) => d.id}
        renderItem={({ item }) => (
          <DoseCard
            dose={item}
            onTake={async () => {
              await takeDose(item.id);
              load();
            }}
            onSkip={async () => {
              await skipDose(item.id);
              load();
            }}
          />
        )}
      />
    </View>
  );
}
