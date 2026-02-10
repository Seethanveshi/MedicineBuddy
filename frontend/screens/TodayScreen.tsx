import { useCallback, useEffect, useState } from "react";
import { View, Text, ScrollView, ActivityIndicator, TouchableOpacity } from "react-native";
import { takeDose, skipDose, getUpcomingDoses, getDosesByDate } from "../api/doses";
import DoseCard from "../components/DoseCard";
import { Dose } from "../types/dose";
import { configureNotifications, cancelDoseNotification, cancelAllDoseNotifications, scheduleUpcomingDoseNotifications,  } from "../utils/notifications";
import { cacheSet } from "@/utils/cache";
import WeekCalendar from "@/components/WeekCalender";
import { format, startOfWeek, addWeeks } from "date-fns";
import {useFocusEffect, useNavigation } from "@react-navigation/native";
import CalendarHeader from "@/components/CalendarHeader";
import { groupByTime } from "@/utils/groupBy";

export default function TodayScreen() {
  const [loading, setLoading] = useState(false);
  const [doses, setDoses] = useState<Dose[]>([]);
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [weekStart, setWeekStart] = useState( startOfWeek(new Date(), { weekStartsOn: 1 }));
  const TODAY_CACHE_KEY = "today_doses";
  const UPCOMING_CACHE_KEY = "upcoming_doses";

 const navigation = useNavigation<any>();


  useEffect(() => {
    load(selectedDate);
  }, [selectedDate]);

  const load = async (date: Date) => {
    try {
      setLoading(true);
      const data = await getDosesByDate(date);
      setDoses(data);
    } catch (e) {
      console.warn("failed to load doses", e);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      load(selectedDate); 
      (async () => {
        const allowed = await configureNotifications();
        if (!allowed) return;

        await cancelAllDoseNotifications();
        const upcoming = await getUpcomingDoses();
        await cacheSet(UPCOMING_CACHE_KEY, upcoming);
        await scheduleUpcomingDoseNotifications(upcoming);
      })();
    }, [selectedDate])
  );

  const groups = groupByTime(doses);

  if (loading) {
    return (
      <View style={{ padding: 16 }}>
        <CalendarHeader
          weekStart={weekStart}
          selectedDate={selectedDate}
          setWeekStart={setWeekStart}
          setSelectedDate={setSelectedDate}
        />
        <WeekCalendar
              selected={selectedDate}
              weekStart={weekStart}
              onSelect={setSelectedDate}
            />
        <ActivityIndicator size="large" style={{ marginTop: 40 }} />
      </View>
    );
  }

  if (!loading && doses.length === 0) {
    return (
      <View style={{ padding: 16 }}>
        <CalendarHeader
          weekStart={weekStart}
          selectedDate={selectedDate}
          setWeekStart={setWeekStart}
          setSelectedDate={setSelectedDate}
        />
        <WeekCalendar
              selected={selectedDate}
              weekStart={weekStart}
              onSelect={setSelectedDate}
        />

        <Text style={{
          textAlign: "center",
          marginTop: 40,
          color: "#666",
          fontSize: 16,
        }}>
          No medicines scheduled 🎉
        </Text>
      </View>
    );
  }

  return (
      <View style={{ flex: 1 }}> 
      <ScrollView >
        <View style={{ padding: 16 }}>
          <CalendarHeader
            weekStart={weekStart}
            selectedDate={selectedDate}
            setWeekStart={setWeekStart}
            setSelectedDate={setSelectedDate}
          />
          <WeekCalendar
              selected={selectedDate}
              weekStart={weekStart}
              onSelect={setSelectedDate}
            />
          <Text style={{ fontSize: 18, marginBottom: 12 }}>
            {format(selectedDate, "MMM dd")}
          </Text>

          {groups.map(([time, doses]) => (
            <View key={time} style={{ marginBottom: 20 }}>
              <Text style={{ fontSize: 18, fontWeight: "600", marginBottom: 10 }}>
                {time}
              </Text>

              {doses.map((dose) => (
                <DoseCard
                  key={dose.id}
                  dose={dose}
                  onTake = {async () => {
                          await takeDose(dose.id);
                          await cancelDoseNotification(dose.id);
                          load(selectedDate);
                      }}
                  onSkip = {async () => {
                          await skipDose(dose.id);
                          await cancelDoseNotification(dose.id);
                      load(selectedDate);
                      }}
                />
              ))}
            </View>
          ))}
        </View>
        </ScrollView>
        <TouchableOpacity
          onPress={() => navigation.navigate("AddMedicine")}
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
            elevation: 5,
          }}
        >
          <Text style={{ color: "white", fontSize: 28 }}>+</Text>
        </TouchableOpacity>
      </View>
  );
}