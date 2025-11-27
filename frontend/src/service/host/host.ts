import { get, writable, type Writable } from "svelte/store";
import { NetService, PacketTypes, type QuestionRevealPacket, type Packet, type HostGamePacket, GameState, type ChangeGameStatePacket, type PlayerJoinPacket, type TickPacket, type QuestionShowPacket, type LeaderboardPacket, LeaderboardEntry, type KickPlayerPacket, PlayerLeavePacket } from "../net";
import type { Player, QuizQuestion } from "../../model/quiz";

export const leaderboard: Writable<LeaderboardEntry[]> = writable([]);
export const state: Writable<GameState> = writable(GameState.Lobby);
export const players: Writable<Player[]> = writable([]);
export const gameCode: Writable<string | null> = writable(null);
export const tick: Writable<number> = writable(0);
export const currentQuestion: Writable<QuizQuestion | null> = writable(null);
export const correctAnswerIndex: Writable<number[]> = writable([]);
export const showAnswer: Writable<boolean> = writable(false);
export const answerCounts: Writable<number[]> = writable([]);
export const isHostNavigating = writable(false);
export const isSigningUp = writable(false);

function resetGameStores() {
    leaderboard.set([]);
    gameCode.set(null);
    players.set([]);
    tick.set(0);
    currentQuestion.set(null);
    correctAnswerIndex.set([]);
    showAnswer.set(false);
    answerCounts.set([]);
}


class HostGame {
    private net: NetService;

    public navigate: ((path: string) => void) | undefined;

    constructor(navigateFunction?: (path: string) => void) {
        this.net = new NetService();
        this.net.onPacket(p => this.onPacket(p));
        if (navigateFunction) {
            this.navigate = navigateFunction;
        }
    }

    hostQuiz(quizId: string) {
        let packet: HostGamePacket = {
            id: PacketTypes.HostGame,
            quizId: quizId,
        }
        this.net.sendPacket(packet);
    }

    connect(token: string) {
        console.log("Host connecting with token...");
        this.net.connect(token);
    }

    public unhost(options: { broadcastEnd?: boolean } = { broadcastEnd: true }) {

        if (options.broadcastEnd) {
            this.net.sendPacket({ id: PacketTypes.HostLeave });
        }
        resetGameStores();
    }

    public signalHostLeaving() {
        this.net.sendPacket({ id: PacketTypes.HostLeave });
        resetGameStores();
    }


    start() {
        this.net.sendPacket({ id: PacketTypes.StartGame });
    }

    public intermission() {
        let packet: ChangeGameStatePacket = {
            id: PacketTypes.ChangeGameState,
            state: GameState.Intermission
        };
        this.net.sendPacket(packet);
    }

    nextQuestion() {
        this.net.sendPacket({ id: PacketTypes.NextQuestion });
    }

    public kickPlayer(playerId: string) {
        let packet: KickPlayerPacket = {
            id: PacketTypes.KickPlayer,
            playerId: playerId
        };
        this.net.sendPacket(packet);
        players.update(p => p.filter(player => player.id !== playerId));
    }


    onPacket(packet: Packet) {
        switch (packet.id) {
            case PacketTypes.HostGame: {
                let data = packet as HostGamePacket;
                gameCode.set(data.quizId);
                break;
            }
            case PacketTypes.ChangeGameState: {
                let data = packet as ChangeGameStatePacket;
                state.set(data.state);
                break;
            }

            case PacketTypes.PlayerJoin: {
                let data = packet as PlayerJoinPacket;
                players.update(p => [...p, data.player]);
                break;
            }

            case PacketTypes.Tick: {
                let data = packet as TickPacket;
                tick.set(data.tick);
                break;
            }

            case PacketTypes.QuestionShow: {
                let data = packet as QuestionShowPacket;
                currentQuestion.set({ ...data.question, index: data.questionIndex });
                break;
            }

            case PacketTypes.QuestionReveal: {
                let data = packet as QuestionRevealPacket;
                correctAnswerIndex.set(data.correctAnswerIndex);
                answerCounts.set(data.answerCounts);
                showAnswer.set(true);
                break;
            }

            case PacketTypes.PlayerLeave: {
                let data = packet as PlayerLeavePacket;

                if (!data || data.playerId === undefined) {
                    break;
                }


                const leavingPlayerId = data.playerId;

                players.update(currentPlayers => {

                    const filteredPlayers = currentPlayers.filter(p => {

                        const idMatch = p.id === leavingPlayerId;

                        return String(p.id) !== String(leavingPlayerId);
                    });
                    return filteredPlayers;
                });

                break;
            }


            case PacketTypes.Leaderboard: {
                let data = packet as LeaderboardPacket;
                leaderboard.set(data.points);
                break;
            }
        }
    }
}
export const game = new HostGame();
    export { GameState, HostGame, resetGameStores };