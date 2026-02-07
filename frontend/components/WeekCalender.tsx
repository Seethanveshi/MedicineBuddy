import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import { addDays, startOfWeek, format, addWeeks } from "date-fns";
import { useState } from "react";

type Props = {
  selected: Date;
  weekStart: Date;
  onSelect: (d: Date) => void;
};

export default function WeekCalendar({ selected,weekStart, onSelect }: Props) {
  const start = weekStart;

  const days = Array.from({ length: 7 }).map((_, i) => addDays(start, i));

  return (
    <View style={styles.row}>
      {days.map((day) => {
        const isSelected =
          format(day, "yyyy-MM-dd") === format(selected, "yyyy-MM-dd");

        return (
          <TouchableOpacity
            key={day.toISOString()}
            onPress={() => onSelect(day)}
            style={[styles.day, isSelected && styles.selected]}
          >
            <Text style={styles.weekday}>{format(day, "EEE")}</Text>
            <Text style={styles.date}>{format(day, "d")}</Text>
          </TouchableOpacity>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 16,
  },
  day: {
    alignItems: "center",
    padding: 8,
    borderRadius: 8,
  },
  selected: {
    backgroundColor: "#1976d2",
  },
  weekday: {
    fontSize: 12,
    color: "#666",
  },
  date: {
    fontSize: 16,
    fontWeight: "600",
  },
});
