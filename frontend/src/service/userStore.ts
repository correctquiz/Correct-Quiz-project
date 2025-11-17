import { get, writable, type Writable } from 'svelte/store';
import { onAuthStateChanged, signOut, type User, type IdTokenResult } from 'firebase/auth';
import { auth } from './firebase';
import { apiService } from './api';

export const isLoggingOut = writable(false);

export interface UserState {
    loggedIn: boolean;
    user: User | null;
    userType: 'Host' | 'Player' | null;
}


export const userStore: Writable<UserState> = writable({
    loggedIn: false,
    user: null,
    userType: null,
});

export const isSigningUp = writable(false);
export const isLoggingIn = writable(false);

onAuthStateChanged(auth, async (firebaseUser) => {
    if (get(isSigningUp) || get(isLoggingIn) || get(isLoggingOut)) {
        return;
    }

    if (firebaseUser) {
        const idTokenResult = await firebaseUser.getIdTokenResult();
        const userRole = idTokenResult.claims.role as 'host' | 'player';
        await apiService.login({ idToken: idTokenResult.token });
        userStore.update(currentStore => {
            return {
                ...currentStore,
                loggedIn: true,
                user: firebaseUser,
                userType: userRole === 'host' ? 'Host' : 'Player',
            };
        });
    } else {
        userStore.set({
            loggedIn: false,
            user: null,
            userType: null,
        });
    }
});

export function logout() {
    signOut(auth);
}