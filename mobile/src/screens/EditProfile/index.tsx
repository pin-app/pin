import React from 'react';
import { View, Text, Image, SafeAreaView, StyleSheet, TouchableOpacity, TextInput } from 'react-native';
import { colors } from '../../theme';
import Feather from '@expo/vector-icons/Feather';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { useNavigation } from '@react-navigation/native';


export default function EditProfileScreen() {
  // doesnt work rn
  const profileImagePlaceholder = 'holder'; 
  const navigation = useNavigation();


  return (
    <SafeAreaView style={styles.container}>
       {/* also don't do anything yet */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Feather name="x" style={styles.headerIcon}/>
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Edit Profile</Text>
        <TouchableOpacity onPress={() => navigation.goBack()}> {/* need to save changes */}
          <FontAwesome6 name="check" style={styles.checkIcon}/>
        </TouchableOpacity>
      </View>

      <View style={styles.profileSection}>
        <Image 
          source={{ uri: profileImagePlaceholder }}
          style={styles.profilePic} 
        />
        <TouchableOpacity>
          <Text style={styles.changePhotoText}>Change Profile Photo</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.inputSection}>
        <View style={styles.inputGroup}>
          <Text style={styles.label}>Name</Text>
          <TextInput 
            placeholder="name"
            style={styles.textInput}
          />
        </View>

        <View style={styles.inputGroup}>
          <Text style={styles.label}>Bio</Text>
          <TextInput 
            placeholder="bio"
            style={styles.textInput}
          />
        </View> 

        <View style={styles.inputGroup}>
          <Text style={styles.label}>Home City</Text>
          <TextInput 
            placeholder="current home city"
            style={styles.textInput}
          />
        </View> 
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background || 'white',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: 10,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border || '#CDCDCD',
  },
  headerTitle: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  headerIcon: {
    fontSize: 35,
  },
  checkIcon: {
    fontSize: 33,
    color: '#3493D9',
  },
  profileSection: {
    alignItems: 'center',
    padding: 20,
  },
  profilePic: {
    width: 150,
    height: 150,
    borderRadius: 75,
    backgroundColor: 'gray',
    marginBottom: 10,
  },
  changePhotoText: {
    padding: 10,
    color: '#3493D9',
  },
  inputSection: {
    paddingHorizontal: 10,
  },
  inputGroup: {
    marginBottom: 20,
  },
  label: {
    paddingBottom: 5,
    fontSize: 16,
    fontWeight: '500',
    color: 'black',
  },
  textInput: {
    fontSize: 16,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border || '#CDCDCD',
    paddingVertical: 5,
  },
});
