<script lang="ts">
  import Router, { location, push } from "svelte-spa-router";
  import { isLoggingOut, userStore } from "./service/userStore";
  import HostQuizListView from "./views/host/HostQuizListView.svelte";
  import HostSignup from "./views/Auth/HostSignup.svelte";
  import HostLogin from "./views/Auth/HostLogin.svelte";
  import PlayerView from "./views/player/PlayerView.svelte";
  import PlayerJoinView from "./views/player/PlayerJoinView.svelte";
  import EditQuizView from "./views/edit/EditQuizView.svelte";
  import PlayerLogin from "./views/Auth/PlayerLogin.svelte";
  import PlayerSignup from "./views/Auth/PlayerSignup.svelte";
  import ForgotPassword from "./views/Auth/ForgotPassword.svelte";
  import { wrap } from "svelte-spa-router/wrap";
  import HostView from "./views/host/HostView.svelte";
  import { tick } from "svelte";
  import VerifyEmailPage from "./views/Auth/VerifyEmailPage.svelte";
    import { get } from "svelte/store";

  $: {
    (async () => {
      await tick();

      if (get(isLoggingOut)) {
        isLoggingOut.set(false); 
        return; 
      }

      const store = $userStore;
      const currentPath = $location;

      if (store.loggedIn) {
        const isOnAuthPage =
          currentPath.includes("login") ||
          currentPath.includes("signup") ||
          currentPath.includes("forgot-password");

        if (isOnAuthPage) {
          if (store.userType === "Host") {
            push("/host");
          } else {
            push("/");
          }
        } else {
          const isPlayer = store.userType === "Player";

          const isHostOnlyPage =
            currentPath.startsWith("/host") || currentPath.startsWith("/edit");

          if (isPlayer && isHostOnlyPage) {
            push("/");
          }
        }
      } else {
        const isAuthProtected =
          currentPath.startsWith("/host") ||
          currentPath.startsWith("/edit") ||
          currentPath.startsWith("/play");

        const isOnAuthPage =
          currentPath.includes("login") ||
          currentPath.includes("signup") ||
          currentPath.includes("forgot-password");

        if (isAuthProtected && !isOnAuthPage) {
          push("/");
        }
      }
    })();
  }

  let routes = {
    "/": PlayerJoinView,
    "/play": PlayerView,
    "/player/login": PlayerLogin,
    "/player/signup": PlayerSignup,
    "/host": HostQuizListView,
    "/host/game": HostView,
    "/host/login": HostLogin,
    "/host/signup": HostSignup,
    "/verify-email": VerifyEmailPage,
    "/edit/:quizId": EditQuizView,
    "/player/forgot-password": wrap({
      component: ForgotPassword,
      props: { userType: "Player" },
    }),
    "/host/forgot-password": wrap({
      component: ForgotPassword,
      props: { userType: "Host" },
    }),
  };
</script>

<Router {routes} />
