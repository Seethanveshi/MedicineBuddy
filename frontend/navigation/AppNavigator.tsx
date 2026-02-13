import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { createNativeStackNavigator } from "@react-navigation/native-stack";

import TodayScreen from "../screens/TodayScreen";
import HistoryScreen from "../screens/HistoryScreen";
import AddMedicineScreen from "../screens/AddMedicineScreen";
import MediTakersScreen from "@/screens/MediTakerScreen";
import AddMediTakerScreen from "@/screens/MediTakerFormScreen";

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
      <Tab.Screen name="Meditaker" component={MediTakersScreen} />
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
      <Stack.Screen
        name="MediTakers"
        component={MediTakersScreen}
        options={{ title: "MediTakers" }}
      />
      <Stack.Screen
        name="MediTakerForm"
        component={AddMediTakerScreen}
        options={{ title: "Add MediTaker" }}
      />

    </Stack.Navigator>
  );
}
