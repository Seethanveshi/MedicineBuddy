import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, ScrollView } from "react-native";
import DateTimePicker from "@react-native-community/datetimepicker";
import { useNavigation } from "@react-navigation/native";
import { createMedicine } from "@/api/doses";


export default function AddMedicineScreen() {
  const [name, setName] = useState("");
  const [dosage, setDosage] = useState("");
  const [time, setTime] = useState("");
  const [showTime, setShowTime] = useState(false);
  const [showStart, setShowStart] = useState(false);
  const [showEnd, setShowEnd] = useState(false);
  const [type, setType] = useState<"everyday" | "custom">("everyday");
  const [startDate, setStartDate] = useState(new Date());
  const [endDate, setEndDate] = useState<Date | null>(null);
  const [weekdays, setWeekdays] = useState<number[]>([]);
  const navigation = useNavigation();

  const validate = () => {
    if (!name.trim()) return "Medicine name required";
    if (!dosage.trim()) return "Dosage required";
    if (type === "custom" && weekdays.length === 0)
        return "Select at least one weekday";
    return null;
  };
  const formatDate = (d: Date) =>
  d.toISOString().split("T")[0];


  return (
    <ScrollView contentContainerStyle={styles.container}>
        <View>
            <Text style={styles.title}>Add Medicine</Text>
            <Text style={styles.label}>Medicine name</Text>
            <TextInput
            value={name}
            onChangeText={setName}
            style={styles.input}
            placeholder="Enter"
            />

            <Text style={styles.label}>Dosage</Text>
            <TextInput
            value={dosage}
            onChangeText={setDosage}
            style={styles.input}
            placeholder="Enter"
            />

            <Text style={styles.label}>Time</Text>

            <TouchableOpacity
            style={styles.input}
            onPress={() => setShowTime(true)}
            >
            <Text>
                {time}
            </Text>
            </TouchableOpacity>

            {showTime && (
            <DateTimePicker
                mode="time"
                value={new Date()}
                onChange={(_, date) => {
                setShowTime(false);
                if (!date) return;

                const hh = date.getHours().toString().padStart(2, "0");
                const mm = date.getMinutes().toString().padStart(2, "0");
                setTime(`${hh}:${mm}`);
                }}
            />
            )}

            <Text style={styles.label}>Schedule</Text>

            <View style={styles.row}>
            <TouchableOpacity onPress={() => setType("everyday")}>
                <Text style={type === "everyday" ? styles.active : styles.inactive}>
                ● Every day
                </Text>
            </TouchableOpacity>

            <TouchableOpacity onPress={() => setType("custom")}>
                <Text style={type === "custom" ? styles.active : styles.inactive}>
                ● Custom
                </Text>
            </TouchableOpacity>
            </View>
            <Text style={styles.label}>Start date</Text>
            <TouchableOpacity style={styles.input} onPress={() => setShowStart(true)}>
            <Text>{startDate.toDateString()}</Text>
            </TouchableOpacity>

            {showStart && (
            <DateTimePicker
                mode="date"
                value={startDate}
                onChange={(_, date) => {
                setShowStart(false);
                if (date) setStartDate(date);
                }}
            />
            )}


            {type === "custom" && (
                <>
                    <Text style={styles.label}>End date</Text>

                    <TouchableOpacity style={styles.input} onPress={() => setShowEnd(true)}>
                    <Text>
                        {endDate ? endDate.toDateString() : "Not set"}
                    </Text>
                    </TouchableOpacity>

                    {showEnd && (
                    <DateTimePicker
                        mode="date"
                        value={endDate || new Date()}
                        onChange={(_, date) => {
                        setShowEnd(false);
                        if (date) setEndDate(date);
                        }}
                    />
                    )}

                    <Text style={styles.label}>Weekdays</Text>

                    <View style={styles.row}>
                    {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map(
                        (day, index) => {
                        const value = index;

                        const selected = weekdays.includes(value);

                        return (
                            <TouchableOpacity
                            key={day}
                            onPress={() => {
                                if (selected) {
                                setWeekdays(weekdays.filter((d) => d !== value));
                                } else {
                                setWeekdays([...weekdays, value]);
                                }
                            }}
                            style={[styles.dayBtn, selected && styles.daySelected]}
                            >
                            <Text>{day}</Text>
                            </TouchableOpacity>
                        );
                        }
                    )}
                    </View>
                </>
            )}

            <TouchableOpacity
                style={styles.save}
                onPress={async () => {
                    const error = validate();
                    if (error) {
                    alert(error);
                    return;
                    }

                    await createMedicine({
                    name,
                    dosage,
                    start_date: formatDate(startDate),
                    end_date: endDate ? formatDate(endDate) : null,
                    schedule: {
                        time,
                        days_of_week:
                        type === "everyday"
                            ? [0,1,2,3,4,5,6]
                            : weekdays,
                    },
                    });

                    navigation.goBack();
            }}
                >
            <Text style={{ color: "white", fontWeight: "600" }}>Save</Text>
            </TouchableOpacity>
        </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: 16,
    paddingBottom: 40,
  },
  title: { fontSize: 22, fontWeight: "700", marginBottom: 20 },

  label: { marginTop: 16, marginBottom: 6, fontWeight: "600" },

  input: {
    borderWidth: 1,
    borderColor: "#ddd",
    padding: 12,
    borderRadius: 8,
  },

  row: {
    flexDirection: "row",
    gap: 12,
    marginTop: 8,
    flexWrap: "wrap",
  },

  active: { color: "#1976d2", fontWeight: "700" },
  inactive: { color: "#777" },

  dayBtn: {
    padding: 8,
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 6,
  },

  daySelected: {
    backgroundColor: "#1976d2",
  },

  save: {
    marginTop: 30,
    backgroundColor: "#1976d2",
    padding: 16,
    alignItems: "center",
    borderRadius: 10,
  },
});
