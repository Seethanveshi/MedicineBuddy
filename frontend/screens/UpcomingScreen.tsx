import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import { API } from "@/api/doses";

export default function UpcomingScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);

  useEffect(() => {
    const fetchDoses = async () => {
        const res = await API.get<Dose[]>("/doses/upcoming");
        setDoses(res.data);
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
