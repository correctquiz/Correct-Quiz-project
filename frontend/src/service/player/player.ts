import { writable, Writable, get } from "svelte/store";
import { NetService, type Packet, PacketTypes, type ConnectPacket, ChangeGameStatePacket, GameState, type QuestionShowPacket, type QuestionAnswerPacket, type PlayerRevealPacket, PlayerJoinPacket, LeaderboardEntry, LeaderboardPacket, type PlayerRankPacket, type PlayerAnswerFeedbackPacket } from "../net";
import type { QuizQuestion } from "../../model/quiz";
import type { Player } from '../../model/quiz';

export const state: Writable<GameState> = writable(GameState.Lobby);
export const points: Writable<number> = writable(0);
export const currentQuestion: Writable<QuizQuestion | null> = writable(null);
export const leaderboard: Writable<LeaderboardEntry[]> = writable([]);
export const rank: Writable<number> = writable(0);
export const maxStreak: Writable<number> = writable(0);
export const streakBonus: Writable<number> = writable(0);
export const currentPlayer: Writable<Player> = writable({
    id: "",
    name: "Loading..."
});

export function resetPlayerStores() {
    state.set(GameState.Lobby);
    points.set(0);
    currentQuestion.set(null);
    leaderboard.set([]);
    rank.set(0);
    maxStreak.set(0);
    streakBonus.set(0);
    currentPlayer.set({ id: "", name: "Loading..." });
}

export class PlayerGame {
    public navigate?: (path: string) => void;
    private messageHandlers: ((packet: Packet) => void)[] = [];
    private net: NetService;

    constructor(navigateFunction?: (path: string) => void) {
        console.trace("Created from here:");
        this.net = new NetService();
        this.net.onPacket(p => this.onPacket(p));
        this.net.onDisconnect((code, reason) => this.handleDisconnect(code, reason));
        if (navigateFunction) {
            this.navigate = navigateFunction;
        }
    }

    private handleDisconnect(code: number, reason: string) {
        if (code !== 1000 && this.navigate) {
            alert("คุณถูกเตะออกจากห้องเกมแล้ว!");
            this.navigate('/');
        }
        resetPlayerStores()
    }

    join(code: string, name: string, token: string) {
        this.net.connect(token);
        let packet: ConnectPacket = {
            id: PacketTypes.Connect,
            code: code,
            name: name,
        }
            this.net.sendPacket(packet);
    }

    answer(questionIndex: number, choiceIndex: number) {
        let packet: QuestionAnswerPacket = {
            id: PacketTypes.Answer,
            question: questionIndex,
            choice: choiceIndex,
        };
        this.net.sendPacket(packet);
    }

    public signalPlayerLeaving() {
        const player = get(currentPlayer);
        if (!player || player.id === "") {
            resetPlayerStores();
            return;
        }

        const leavePacket = {
            id: PacketTypes.PlayerLeave,
            playerId: player.id
        };

        this.net.sendPacket(leavePacket);

        resetPlayerStores();
    }

    onMessage(handler: (packet: Packet) => void) {
        this.messageHandlers.push(handler);
    }


    onPacket(packet: Packet) {
        switch (packet.id) {
            case PacketTypes.PlayerJoin: {
                console.log("DEBUG: PlayerJoin Packet Received!");
                let data = packet as PlayerJoinPacket;
                console.log("DEBUG: Player Data:", data.player);
                currentPlayer.set(data.player);
                if (this.navigate) {
                    this.navigate('/play');
                }
                break;
            }
            case PacketTypes.ChangeGameState: {
                let data = packet as ChangeGameStatePacket;
                if (data.state === GameState.GameEndedState) {
                    alert(`เกมถูกปิดโดย Host แล้ว`);
                    if (this.navigate) {
                        this.navigate('/');
                    }
                }
                state.set(data.state);
                break;
            }
            case PacketTypes.QuestionShow: {
                let data = packet as QuestionShowPacket;
                currentQuestion.set({ ...data.question, index: data.questionIndex });
                break;
            }
            case PacketTypes.PlayerAnswerFeedback: {
                this.messageHandlers.forEach(handler => handler(packet));
                let data = packet as PlayerAnswerFeedbackPacket;
                maxStreak.set(data.maxStreak);
                streakBonus.set(data.streakBonus);
                break;
            }
            case PacketTypes.PlayerReveal: {
                let data = packet as PlayerRevealPacket;
                points.set(data.points);
                break;
            }
            case PacketTypes.Leaderboard: {
                let data = packet as LeaderboardPacket;
                leaderboard.set(data.points);

                const myPlayer = get(currentPlayer);
                const myRank = data.points.findIndex(p => p.name === myPlayer?.name);

                if (myRank !== -1) {
                    rank.set(myRank + 1);
                }
                break;
            }
            case PacketTypes.PlayerRank: {
                let data = packet as PlayerRankPacket;
                rank.set(data.rank);
                break;
            }
        }
    }
}

export const game = new PlayerGame();
