import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { Dose } from "../types/dose";
import DoseCard from "../components/DoseCard";
import axios from "axios";
import { API } from "@/api/doses";

export default function HistoryScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);

  useEffect(() => {
    const fetchHistoryDoses = async () => {
        const res = await API.get<Dose[]>("/doses/history");
        setDoses(res.data);
    };
    
    fetchHistoryDoses();
  }, []);

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 12 }}>History</Text>

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
