import { View, Text, Button, StyleSheet, TouchableOpacity } from "react-native";
import { Dose } from "../types/dose";
import { useNavigation } from "@react-navigation/native";

type Props = {
  dose: Dose;
  onTake: () => void;
  onSkip: () => void;
};

export default function DoseCard({ dose, onTake, onSkip }: Props) {
  const style = STATUS_STYLES[dose.status];
  const navigation = useNavigation<any>()
  return (
    <View style={[
        styles.card,
        {
          borderLeftColor: style.border,
          backgroundColor: style.bg,
        },
      ]}>
       <TouchableOpacity
          onPress={() => navigation.navigate("AddMedicine", {id: dose.medicine_id})}
        >
        <Text style={styles.name}>
          {dose.name}
        </Text>
        <Text style={styles.dosage}>
          {dose.dosage}
        </Text>

        <Text style={{ color: style.text, marginTop: 4 }}>
          {dose.status.toUpperCase()}
        </Text>

        {dose.status === "pending" && (
          <View style={styles.actions}>
            <Button title="Take" onPress={onTake} />
            <Button title="Skip" onPress={onSkip} />
          </View>
        )}
      </TouchableOpacity>
    </View>
  );
}

// function formatTimeThenDate(datetime: string) {
//   const date = new Date(datetime);
//   const time = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
//   const day = date.toLocaleDateString([], { day: '2-digit', month: 'short', year: 'numeric' });
//   return `${time} : ${day}`;
// }

const styles = StyleSheet.create({
  card: {
    padding: 16,
    marginBottom: 12,
    borderRadius: 12,
    borderLeftWidth: 5,
  },
  name: {
    fontSize: 18,
    fontWeight: "700",
    color: "#333",
  },
  dosage: {
    fontSize: 16,
    fontWeight: "500",
    color: "#555",
    marginTop: 2,
  },
  time: {
    fontSize: 20,
    fontWeight: "600",
  },
  status: {
    marginTop: 6,
    color: "#666",
  },
  actions: {
    marginTop: 10,
    flexDirection: "row",
    justifyContent: "space-between",
  },
});

const STATUS_STYLES = {
  pending: {
    border: "#f9a825",
    bg: "#fff8e1",
    text: "#f57f17",
  },
  taken: {
    border: "#2e7d32",
    bg: "#e8f5e9",
    text: "#1b5e20",
  },
  missed: {
    border: "#c62828",
    bg: "#ffebee",
    text: "#b71c1c",
  },
  skipped: {
    border: "#9e9e9e",
    bg: "#f5f5f5",
    text: "#616161",
  },
};


