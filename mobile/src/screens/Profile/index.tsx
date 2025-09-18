import React from 'react';
import { View, StyleSheet, Text, Image } from 'react-native';
import { colors, spacing } from '../../theme';
import { PageTitle, Button } from '../../shared/components';

export default function ProfileScreen() {
  const handleEditProfile = () => {
    // useless for now
    console.log('Edit profile pressed');
  };


  return (
    <View style={styles.container}>
      {/* <View style={styles.topSection}></View> */}
      <Image style={styles.profilePicWrapper}/>
      <View style={styles.bottomSection}>
        <View style={styles.profileCard}>
        <View style={styles.statsRow}>
          <View style={styles.statsCol}>
            <Text style={styles.statValue}>5</Text>
            <Text style={styles.statLabel}>Places</Text>
          </View>

          <View style={styles.statsCol}>
            <Text style={styles.statValue}>1.7k</Text>
            <Text style={styles.statLabel}>Following</Text>
          </View>

          <View style={styles.statsCol}>
            <Text style={styles.statValue}>2.3k</Text>
            <Text style={styles.statLabel}>Followers</Text>
          </View>
        </View>
      </View>
      </View>
      <Button 
        title="Edit Profile" 
        onPress={handleEditProfile}
        variant="primary"
        size="md"
        style={styles.button}
      />
      <Button 
        title="Share" 
        onPress={handleEditProfile}
        variant="secondary"
        size="md"
        style={styles.button}
      />
      <Button 
        title="Share but with a border around it" 
        onPress={handleEditProfile}
        variant="outline"
        size="md"
        style={styles.button}
      />
    </View>
  );
}



const styles = StyleSheet.create({
  container: {
    display: 'flex', //grow and stretch
    alignItems: 'center', // aligns items in the middle of the x y cross
    justifyContent: 'center',
    flex: 1,
    position: 'relative',
    
  },
  
  // topSection: {
  //   flex: 0.75,
  //   backgroundColor: 'darkolivegreen'
  // },

  bottomSection: {
    flex: 1,
  },

  profileCard: {
    // top: 170,
    alignSelf: 'center',
    width: 330,
    height: 100,
    backgroundColor: 'white',
    borderRadius: 25,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: 'rgba(0,0,0,0.08)',
    shadowColor: '#000',
    shadowOpacity: 0.12,
    shadowRadius: 12,
    shadowOffset: { width: 0, height: 6 },
    overflow: 'hidden',
  },

  statsRow: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-around',
    paddingHorizontal: 16,
  },

  statsCol: {
    alignItems: 'center',
  },

  statLabel: {
    marginTop: 2,
    fontSize: 12,
    color: '#6B7280',
  },


  statValue: {
    fontSize: 20,
    fontWeight: '700',
    color: '#111827'
  },

  divider: {
    width: StyleSheet.hairlineWidth,
    height: '60%',
    backgroundColor: '#E5E7B',
  },

  profilePicWrapper: {
    width: 150,
    height: 150,
    borderRadius: '100%',
    overflow: 'hidden',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'gray',
  },
  button: {
    marginTop: spacing.lg,
  },

  profilePic: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
  }
});
