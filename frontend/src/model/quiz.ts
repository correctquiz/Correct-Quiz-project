export interface Quiz {
    id: number;
    name: string;
    questions: QuizQuestion[];
}

export interface User {
    id: number;
    email: string;
    role: 'host' | 'player';
}

export interface Player {
    id: string;
    name: string;
}

export interface QuizQuestion {
    index: any;
    id: number;
    name: string;
    time: number;
    choices: QuizChoice[];
    correctAnswerIndex: number;
    imageUrl?: string;
}

export interface QuizChoice {
    id: number;
    name: string;
    correct: boolean;
    imageUrl?: string;
}

export interface RegisterError {
    errorStatus: number;
}

export const COLORS = ["bg-pink-400", "bg-orange-200", "bg-green-200", "bg-purple-200"];