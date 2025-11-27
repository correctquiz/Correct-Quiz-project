import { Player, QuizQuestion } from "../model/quiz";

export enum PacketTypes {
    Connect = 0,
    HostGame = 1,
    StartGame = 5,
    Answer = 7,
    AnswerReceived = 9,
    QuestionShow = 2,
    ChangeGameState = 3,
    PlayerJoin = 4,
    Tick = 6,
    PlayerAnswerFeedback = 8,
    QuestionReveal = 10,
    PlayerReveal = 11,
    Leaderboard = 12,
    NextQuestion = 13,
    PlayerRank = 14,
    KickPlayer = 15,
    HostLeave = 16,
    PlayerLeave = 17,
}

export enum GameState {
    Lobby,
    Play,
    Intermission,
    Reveal,
    End,
    GameEndedState,
}

export interface Packet {
    id: PacketTypes;
}

export interface KickPlayerPacket extends Packet {
    playerId: string;
}

export interface HostGamePacket extends Packet {
    quizId: string;
}

export interface ChangeGameStatePacket extends Packet {
    state: GameState;
    code?: string;
}
export interface PlayerJoinPacket extends Packet {
    player: Player;
}

export interface TickPacket extends Packet {
    tick: number;
}

export interface ConnectPacket extends Packet {
    code: string;
    name: string;
}

export interface QuestionShowPacket extends Packet {
    question: QuizQuestion;
    questionIndex: number;
}

export interface QuestionAnswerPacket extends Packet {
    question: number;
    choice: number;
}

export interface PlayerRevealPacket extends Packet {
    points: number;
}

export interface PlayerAnswerFeedbackPacket extends Packet {
    correctAnswerIndex: number[];
    isCorrect: boolean;
    streakBonus: number;
    maxStreak: number;
}

export interface QuestionRevealPacket extends Packet {
    question: QuizQuestion;
    correctAnswerIndex: number[];
    answerCounts: number[];
    maxStreak: number;
}

export interface LeaderboardEntry {
    name: string;
    points: number;
    correctCount: number;
}

export interface LeaderboardPacket extends Packet {
    points: LeaderboardEntry[];
}

export interface PlayerRankPacket extends Packet {
    rank: number;
}

export interface NextQuestionPacket extends Packet { }

export interface PlayerLeavePacket extends Packet {
    playerId: string;
}


export class NetService {
    private webSocket!: WebSocket;
    private textDecoder: TextDecoder = new TextDecoder();
    private textEncoder: TextEncoder = new TextEncoder();

    private onPacketCallback?: (packet: any) => void;
    public onDisconnectCallback?: (code: number, reason: string) => void

    private pendingQueue: Uint8Array[] = [];

    constructor() {
        console.log("DEBUG: NetService initialized. Waiting for connect() call.");
    }
    

    connect(token: string) {
        let baseUrl = (import.meta as any).env.VITE_WS_URL || "ws://localhost:3000/ws";
        const WS_URL = `${baseUrl}?token=${token}`;
        console.log("📢 WS URL is:", WS_URL);
        console.log("📢 Connecting WS with token:", WS_URL);
        this.webSocket = new WebSocket(WS_URL);
        this.webSocket.onopen = () => {
            console.log("opened connection");
            while (this.pendingQueue.length > 0) {
                const data = this.pendingQueue.shift();
                if (data) {
                    console.log("🚀 Sending queued packet...");
                    this.webSocket.send(data);
                }
            }
        };

        this.webSocket.onmessage = async (event: MessageEvent) => {
            const arrayBuffer = await event.data.arrayBuffer();
            const bytes = new Uint8Array(arrayBuffer);
            const packetId = bytes[0];

            const packet = JSON.parse(this.textDecoder.decode(bytes.subarray(1)));

            packet.id = packetId;

            if (this.onPacketCallback)
                this.onPacketCallback(packet);
        }

        this.webSocket.onclose = (event) => {
            if (this.onDisconnectCallback) {
                this.onDisconnectCallback(event.code, event.reason);
            }
        };
    }

    public onDisconnect(callback: (code: number, reason: string) => void) {
        this.onDisconnectCallback = callback;
    }

    public disconnect() {
        if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
            this.webSocket.close();
            console.log("WebSocket connection closed.");
        }
    }



    onPacket(callback: (packet: Packet) => void) {
        this.onPacketCallback = callback;
    }

    sendPacket(packet: Packet) {
        if (!this.webSocket || this.webSocket.readyState !== WebSocket.OPEN) {
            console.warn("⚠️ Socket not ready. Skipping packet:", packet.id);
            return;
        }
        const packetId = packet.id;
        const packetData = JSON.stringify(packet, (key, value) =>
            key == "id" ? undefined : value
        );

        const packetIdArray = new Uint8Array([packetId]);
        const packetDataArray = this.textEncoder.encode(packetData);

        const mergedArray = new Uint8Array(
            packetIdArray.length + packetDataArray.length,
        );
        mergedArray.set(packetIdArray);
        mergedArray.set(packetDataArray, packetIdArray.length);

        if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
            this.webSocket.send(mergedArray);
        } else {
            console.log("⏳ Socket not ready. Queuing packet:", packet.id);
            this.pendingQueue.push(mergedArray);
        }

        this.webSocket.send(mergedArray);
    }
}