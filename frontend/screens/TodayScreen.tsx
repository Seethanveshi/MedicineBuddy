import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { takeDose, skipDose, getUpcomingDoses, getDosesByDate } from "../api/doses";
import DoseCard from "../components/DoseCard";
import { Dose } from "../types/dose";
import { configureNotifications, scheduleDoseNotification, cancelDoseNotification, cancelAllDoseNotifications, scheduleUpcomingDoseNotifications,  } from "../utils/notifications";
import { cacheGet, cacheSet } from "@/utils/cache";
import WeekCalendar from "@/components/WeekCalender";
import { format } from "date-fns";
import { ActivityIndicator } from "react-native";
import { Gesture, GestureDetector } from "react-native-gesture-handler";
import { startOfWeek, addWeeks } from "date-fns";

export default function TodayScreen() {
  const [loading, setLoading] = useState(false);

  const [doses, setDoses] = useState<Dose[]>([]);
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [weekStart, setWeekStart] = useState( startOfWeek(new Date(), { weekStartsOn: 1 }));
  const TODAY_CACHE_KEY = "today_doses";
  const UPCOMING_CACHE_KEY = "upcoming_doses";

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

  useEffect(() => {
    const init = async () => {
      // 1️⃣ Load cached today doses
      const cached = await cacheGet<Dose[]>(TODAY_CACHE_KEY);
      if (cached) setDoses(cached);

      // 2️⃣ Fetch fresh today doses
      try {
        const upcomingData = await getUpcomingDoses();
        await cacheSet(UPCOMING_CACHE_KEY, upcomingData);

        // 3️⃣ Schedule notifications for today's pending doses
        // await Promise.all(
        //   todayData.filter(d => d.status === "pending").map(scheduleDoseNotification)
        // );

        await cancelAllDoseNotifications();
        await scheduleUpcomingDoseNotifications(upcomingData);
      } catch (e) {
        console.warn("Failed to fetch today doses, using cache", e);
      }

      // 4️⃣ Ensure upcoming notifications are scheduled
      try {
        const allowed = await configureNotifications();
        if (!allowed) return;

        const upcomingData = await getUpcomingDoses();
        await cacheSet(UPCOMING_CACHE_KEY, upcomingData);
        await scheduleUpcomingDoseNotifications(upcomingData);
      } catch (e) {
        console.warn("Failed to schedule upcoming notifications", e);
      }
    };

    init();
  }, []);

  useEffect(() => {
    const initNotifications = async () => {
      const allowed = await configureNotifications();
      if (!allowed) return;

      const upcoming = await getUpcomingDoses();

      await cancelAllDoseNotifications();
      await scheduleUpcomingDoseNotifications(upcoming);
    };

    initNotifications();
  }, []);


  const groupByTime = (doses: Dose[]) => {
    const map: Record<string, Dose[]> = {};

    doses.forEach((dose) => {
      const time = new Date(dose.scheduled_at).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });

      if (!map[time]) map[time] = [];
        map[time].push(dose);
    });

    return Object.entries(map).sort(([a], [b]) => {
      return a.localeCompare(b);
    });
  };

  const groups = groupByTime(doses);

  // const swipeGesture = Gesture.Pan()
  // .onEnd((event) => {
  //   if (event.translationX > 50) {
  //     // swipe right → previous day
  //     setSelectedDate((prev) => addDays(prev, -1));
  //   }

  //   if (event.translationX < -50) {
  //     // swipe left → next day
  //     setSelectedDate((prev) => addDays(prev, 1));
  //   }
  // });

  if (loading) {
    return (
      <View style={{ padding: 16 }}>
        <View style={{
              flexDirection: "row",
              justifyContent: "space-between",
              alignItems: "center",
              marginBottom: 10,
            }}>
          <Text onPress={() => {
            const prev = addWeeks(weekStart, -1);
            setWeekStart(prev);
            setSelectedDate(prev);  
          }}>
            ◀
          </Text>

          <Text style={{ fontSize: 16, fontWeight: "600" }}>
            {format(selectedDate, "MMM yyyy")}
          </Text>

          <Text onPress={() => {
            const next = addWeeks(weekStart, 1);
            setWeekStart(next);
            setSelectedDate(next);   
          }}>
            ▶
          </Text>
        </View>
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
        <View style={{
              flexDirection: "row",
              justifyContent: "space-between",
              alignItems: "center",
              marginBottom: 10,
            }}>
          <Text onPress={() => {
            const prev = addWeeks(weekStart, -1);
            setWeekStart(prev);
            setSelectedDate(prev);
          }}>
            ◀
          </Text>

          <Text style={{ fontSize: 16, fontWeight: "600" }}>
            {format(selectedDate, "MMM yyyy")}
          </Text>

          <Text onPress={() => {
            const next = addWeeks(weekStart, 1);
            setWeekStart(next);
            setSelectedDate(next); 
          }}>
            ▶
          </Text>
        </View>
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
        <View style={{ padding: 16 }}>
          <View style={{
              flexDirection: "row",
              justifyContent: "space-between",
              alignItems: "center",
              marginBottom: 10,
            }}>
              <Text onPress={() => {
                const prev = addWeeks(weekStart, -1);
                setWeekStart(prev);
                setSelectedDate(prev);   // 👈 IMPORTANT
              }}>
                ◀
              </Text>

              <Text style={{ fontSize: 16, fontWeight: "600" }}>
                {format(selectedDate, "MMM yyyy")}
              </Text>

              <Text onPress={() => {
                const next = addWeeks(weekStart, 1);
                setWeekStart(next);
                setSelectedDate(next);   // 👈 IMPORTANT
              }}>
                ▶
              </Text>
            </View>
          <WeekCalendar
              selected={selectedDate}
              weekStart={weekStart}
              onSelect={setSelectedDate}
            />

          <Text style={{ fontSize: 18, marginBottom: 12 }}>
            {format(selectedDate, "MMM d")}
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
      </View>
  );
}

