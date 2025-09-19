import React from 'react';
import { View, StyleSheet, Text, SafeAreaView, TouchableOpacity } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing } from '@/theme';
import Button from '@/components/Button';
import ProfileHeader from './components/ProfileHeader';

export default function ProfileScreen() {
  const handleEditProfile = () => {
    console.log('Edit profile pressed');
  };

  const handleMenuPress = () => {
    console.log('Menu pressed - will eventually open settings');
  };

  return (
    <SafeAreaView style={styles.container}>
      <ProfileHeader username="raquentin" onMenuPress={handleMenuPress} />
      
      <View style={styles.content}>
        <View style={styles.profilePicWrapper}/>
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
        
        <View style={styles.buttonRow}>
          <Button 
            title="Edit Profile" 
            onPress={handleEditProfile}
            variant="primary"
            size="md"
            style={styles.editButton}
          />
          <Button 
            title="Share Profile" 
            onPress={handleEditProfile}
            variant="secondary"
            size="md"
            style={styles.shareButton}
          />
          <TouchableOpacity style={styles.followButton} onPress={handleEditProfile}>
            <FontAwesome6 name="user-plus" size={20} color={colors.text} />
          </TouchableOpacity>
        </View>
      </View>
    </SafeAreaView>
  );
}



const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    flex: 1,
    alignItems: 'center',
    paddingHorizontal: spacing.md,
  },

  bottomSection: {
    marginTop: spacing.lg,
  },

  profileCard: {
    alignSelf: 'center',
    width: 330,
    height: 80,
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
    paddingHorizontal: 8,
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
  buttonRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: spacing.md,
    paddingHorizontal: spacing.md,
    gap: spacing.sm,
  },
  editButton: {
    flex: 1,
  },
  shareButton: {
    flex: 1,
  },
  followButton: {
    width: 44,
    height: 44,
    borderRadius: 8,
    backgroundColor: colors.background,
    borderWidth: 1,
    borderColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
  },

  profilePic: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
  }
});
