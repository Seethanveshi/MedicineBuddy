import AppNavigator from '@/navigation/AppNavigator';
import 'react-native-reanimated';
import { GestureHandlerRootView } from "react-native-gesture-handler";

export const unstable_settings = {
  anchor: '(tabs)',
};

export default function RootLayout() {

  return (
    // <GestureHandlerRootView style={{ flex: 1 }}>
      <AppNavigator />
    /* </GestureHandlerRootView> */
  );
}
