import type { Quiz, User } from "../model/quiz";

interface RegisterError {
    errorStatus: number;
}

export const BASE_URL = "https://correct-quiz-project.onrender.com";

export class ApiService {


    async getQuizById(id: number): Promise<Quiz | null> {
        let response = await fetch(`${BASE_URL}/api/quizzes/${id}`, {
            credentials: 'include'
        });
        if (!response.ok) {
            return null;
        }

        let json = await response.json();
        return json;
    }


    async getQuizzes(): Promise<Quiz[]> {
        let response = await fetch('${BASE_URL}/api/quizzes', {
            credentials: 'include'
        });
        if (!response.ok) {
            alert("Failed to fetch quizzes!");
            return [];
        }

        let json = await response.json();
        return json;
    }

    async saveQuiz(quizId: number, quiz: Quiz) {
        let response = await fetch(`${BASE_URL}/api/quizzes/${quizId}`, {
            method: "PUT",
            body: JSON.stringify(quiz),
            headers: {
                "Content-Type": "application/json"
            }, credentials: 'include'
        });

        if (!response.ok) {
            alert("Save ชุดคำถามไม่สำเร็จ")
        }
        return response.ok;
    }

    async login(credentials: any): Promise<{ token: string } | null> {
        const response = await fetch(`${BASE_URL}/api/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials),
        });

        if (!response.ok) {
            alert("Login ไม่สำเร็จ!");
            return null;
        }
        return await response.json();
    }

    async getEmailForUsername(username: string): Promise<string | null> {
        const response = await fetch(`${BASE_URL}/api/users/email/${username}`);

        if (!response.ok) {
            if (response.status === 404) {
                console.log(`ไม่พบผู้ใช้ที่มี username '${username}'`);
            } else {
                console.error("การดึงอีเมลสำหรับ username ล้มเหลว");
            }
            return null;
        }

        const data = await response.json();
        return data.email;
    }

    async register(userData: { email: string, username: string, firebaseUid: string, role: string }): Promise<User | null> {
        const response = await fetch(`${BASE_URL}/api/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            if (response.status === 409) {
                const errorJson = await response.json();

                const errorData = {
                    response: {
                        status: response.status, 
                        data: {
                            error: errorJson.error || "Unknown backend error"
                        }
                    }
                };
                throw errorData; 
            }
            return await response.json();
        }
        return await response.json();
    }

    async verifyEmailToken(token: string): Promise < any > {
    const response = await fetch(`${BASE_URL}/api/auth/verify-email`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token: token }),
    });

    if(!response.ok) {
    const errorJson = await response.json();
    const errorData = {
        response: {
            data: {
                error: errorJson.error || "Unknown verification error"
            }
        }
    };
    throw errorData;
}

return await response.json(); 
    }

    
    async setInitialClaims(data: { firebaseUid: string; username: string; role: string }): Promise < void> {
    const response = await fetch(`${BASE_URL}/set-initial-claims`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if(!response.ok) {
    const errorData = { response: { status: response.status } };
    throw errorData;
}
    }
    
async resendVerificationEmail(email: string): Promise < any > {
    const response = await fetch(`${BASE_URL}/api/auth/resend-verification`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email }),
    });

    if(!response.ok) {
    const errorJson = await response.json();
    const errorData = {
        response: {
            data: {
                error: errorJson.error || "Unknown resend error"
            }
        }
    };
    throw errorData;
}
return await response.json();
    }

    async deleteQuestion(id: number): Promise < boolean > {
    const response = await fetch(`${BASE_URL}/api/questions/${id}`, {
        method: 'DELETE',
        credentials: 'include'
    });
    return response.ok;
}

    async createQuiz(name: string): Promise < Quiz | null > {
    const response = await fetch(`${BASE_URL}/api/quizzes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: name }),
        credentials: 'include',
    });
    if(!response.ok) {
    if (response.status === 401) {
        alert("เซสชันของคุณหมดอายุ กรุณาล็อกอินใหม่อีกครั้ง");
    } else {
        alert("สร้างควิซไม่สำเร็จ!");
    }
    return null;
}
return await response.json();
    }

    async deleteQuiz(id: number): Promise < boolean > {
    const response = await fetch(`${BASE_URL}/api/quizzes/${id}`, {
        method: 'DELETE', credentials: 'include'
    });
    return response.ok;
}

}
export const apiService = new ApiService();