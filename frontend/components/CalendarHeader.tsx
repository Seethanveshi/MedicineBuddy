import { View, Text } from "react-native";
import { addWeeks, format } from "date-fns";

type Props = {
  weekStart: Date;
  selectedDate: Date;
  setWeekStart: (d: Date) => void;
  setSelectedDate: (d: Date) => void;
};

export default function CalendarHeader({
  weekStart,
  selectedDate,
  setWeekStart,
  setSelectedDate,
}: Props) {
  return (
    <View
      style={{
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        marginBottom: 10,
      }}
    >
      <Text
        onPress={() => {
          const prev = addWeeks(weekStart, -1);
          setWeekStart(prev);
          setSelectedDate(prev);
        }}
      >
        ◀
      </Text>

      <Text style={{ fontSize: 16, fontWeight: "600" }}>
        {format(selectedDate, "MMM yyyy")}
      </Text>

      <Text
        onPress={() => {
          const next = addWeeks(weekStart, 1);
          setWeekStart(next);
          setSelectedDate(next);
        }}
      >
        ▶
      </Text>
    </View>
  );
}
