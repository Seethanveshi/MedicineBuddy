import AsyncStorage from "@react-native-async-storage/async-storage";

export const cacheSet = async (key: string, value: any) => {
  try {
    await AsyncStorage.setItem(key, JSON.stringify(value));
  } catch (e) {
    console.warn("cache set failed", e);
  }
};

export const cacheGet = async <T>(key: string): Promise<T | null> => {
  try {
    const raw = await AsyncStorage.getItem(key);
    return raw ? JSON.parse(raw) : null;
  } catch (e) {
    console.warn("cache get failed", e);
    return null;
  }
};
