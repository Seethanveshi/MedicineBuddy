import { useCallback, useEffect, useState } from "react";
import { View, Text, ActivityIndicator, TouchableOpacity, Alert } from "react-native";
import { deleteMediTaker, listMediTakers } from "@/api/doses";
import { MediTaker } from "@/types/meditaker";
import { useFocusEffect, useNavigation } from "@react-navigation/native";

export default function MediTakersScreen() {
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<MediTaker[]>([]); 
  const navigation = useNavigation<any>()

  const load = async () => {
    try {
      setLoading(true);
      const data = await listMediTakers();
      setItems(data);
    } catch (e) {
      console.warn("failed to load medi-takers", e);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      load();
    }, [])
  );

  const onDelete = (id: string, name: string) => {
    Alert.alert(
      "Remove MediTaker",
      `Remove ${name}?`,
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Delete",
          style: "destructive",
          onPress: async () => {
            await deleteMediTaker(id);
            load(); // refresh list
          },
        },
      ]
    );
  };

  
  if (loading) {
    return <ActivityIndicator style={{ marginTop: 40 }} />;
  }


  if (!loading && items.length === 0) {
    return (
      <View style={{ padding: 16 }}>
        <Text>No medi-takers added yet</Text>
      </View>
    );
  }


  return (
    <View style={{ padding: 16 }}>
      {items.map((m) => (
        <View
          key={m.id}
          style={{
            padding: 12,
            borderWidth: 1,
            borderColor: "#eee",
            borderRadius: 8,
            marginBottom: 12,
          }}
        >
          <Text style={{ fontWeight: "700" }}>{m.name}</Text>
          <Text>{m.relationship}</Text>
          <Text style={{ color: "#666" }}>{m.email}</Text>

          <TouchableOpacity
            onPress={() => onDelete(m.id, m.name)}
            style={{ marginTop: 10 }}
          >
            <Text style={{ color: "red" }}>Delete</Text>
          </TouchableOpacity>
        </View>
      ))}

      <TouchableOpacity
          onPress={() => navigation.navigate("AddMediTaker")}
          style={{
            position: "absolute",
            bottom: 20,
            right: 20,
            width: 60,
            height: 60,
            borderRadius: 30,
            backgroundColor: "#1976d2",
            justifyContent: "center",
            alignItems: "center",
          }}
      >
          <Text style={{ color: "white", fontSize: 28 }}>+</Text>
      </TouchableOpacity>
    </View>
  );
}