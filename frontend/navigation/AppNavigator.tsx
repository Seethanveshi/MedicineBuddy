import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { createNativeStackNavigator } from "@react-navigation/native-stack";

import TodayScreen from "../screens/TodayScreen";
import UpcomingScreen from "../screens/UpcomingScreen";
import HistoryScreen from "../screens/HistoryScreen";
import AddMedicineScreen from "../screens/AddMedicineScreen";

const Tab = createBottomTabNavigator();
const Stack = createNativeStackNavigator();

export type RootStackParamList = {
  Main: undefined;
  AddMedicine: undefined;
};

function Tabs() {
  return (
    <Tab.Navigator>
      <Tab.Screen name="Today" component={TodayScreen} />
      <Tab.Screen name="Upcoming" component={UpcomingScreen} />
      <Tab.Screen name="History" component={HistoryScreen} />
    </Tab.Navigator>
  );
}

export default function AppNavigator() {
  return (
    <Stack.Navigator>
      {/* Tabs are main */}
      <Stack.Screen
        name="Main"
        component={Tabs}
        options={{ headerShown: false }}
      />

      {/* Add Medicine page */}
      <Stack.Screen
        name="AddMedicine"
        component={AddMedicineScreen}
        options={{ title: "Add Medicine" }}
      />
    </Stack.Navigator>
  );
}
