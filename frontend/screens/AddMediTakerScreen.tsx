import { useState } from "react";
import { View, Text, TextInput, TouchableOpacity } from "react-native";
import { useNavigation } from "@react-navigation/native";
import { createMediTaker } from "@/api/doses";
import { StyleSheet } from "react-native";


export default function AddMediTakerScreen() {
  const navigation = useNavigation<any>();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [relationship, setRelationship] = useState("");

  return (
    <View style={{ padding: 16 }}>
      <Text style={{ fontSize: 20, fontWeight: "700", marginBottom: 20 }}>
        Add MediTaker
      </Text>

      <Text>Name</Text>
        <TextInput
            value={name}
            onChangeText={setName}
            placeholder="Enter"
            style={styles.input}
        />

        <Text>Email</Text>
        <TextInput
            value={email}
            onChangeText={setEmail}
            placeholder="Enter"
            autoCapitalize="none"
            style={styles.input}
        />

        <Text>Relationship</Text>
        <TextInput
            value={relationship}
            onChangeText={setRelationship}
            placeholder="Enter"
            style={styles.input}
        />

        <TouchableOpacity
            style={styles.save}
            onPress={async () => {
                if (!name || !email) {
                alert("Name & email required");
                return;
                }

                await createMediTaker({ name, email, relationship });

                navigation.goBack();
            }}
            >
            <Text style={{ color: "white" }}>Save</Text>
        </TouchableOpacity>

    </View>
  );
}



const styles = StyleSheet.create({
  input: {
    borderWidth: 1,
    borderColor: "#ddd",
    padding: 12,
    borderRadius: 8,
    marginBottom: 16,
  },
  save: {
    marginTop: 20,
    backgroundColor: "#1976d2",
    padding: 16,
    alignItems: "center",
    borderRadius: 10,
  },
});