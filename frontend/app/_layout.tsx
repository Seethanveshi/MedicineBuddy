import AppNavigator from '@/navigation/AppNavigator';
import { NavigationContainer } from '@react-navigation/native';
import 'react-native-reanimated';


export const unstable_settings = {
  anchor: '(tabs)',
};

export default function RootLayout() {

  return (
      <AppNavigator />
  );
}
