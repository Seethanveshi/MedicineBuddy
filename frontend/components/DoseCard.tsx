import { View, Text, Button, StyleSheet } from "react-native";
import { Dose } from "../types/dose";

type Props = {
  dose: Dose;
  onTake: () => void;
  onSkip: () => void;
};

export default function DoseCard({ dose, onTake, onSkip }: Props) {
  return (
    <View style={styles.card}>
      <Text style={styles.time}>
        {formatTimeThenDate(dose.scheduled_at)}
      </Text>

      <Text style={styles.status}>{dose.status}</Text>

      {dose.status === "pending" && (
        <View style={styles.actions}>
          <Button title="Take" onPress={onTake} />
          <Button title="Skip" onPress={onSkip} />
        </View>
      )}
    </View>
  );
}

function formatTimeThenDate(datetime: string) {
  const date = new Date(datetime);
  const time = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  const day = date.toLocaleDateString([], { day: '2-digit', month: 'short', year: 'numeric' });
  return `${time} : ${day}`;
}

const styles = StyleSheet.create({
  card: {
    padding: 16,
    marginBottom: 12,
    borderRadius: 10,
    backgroundColor: "#fff",
    elevation: 2,
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
