"use client";

import { useState, useEffect } from "react";
import Header from "../../components/Header";
import Link from "next/link";
import Image from "next/image";
import { QuizService, formatDate } from "../lib/quiz-service";
import { listUserAttempts } from "../lib/quizApi";

export default function HistoryPage() {
  const [attempts, setAttempts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadAttempts = async () => {
      try {
        // Load actual quiz attempts from backend
        const data = await listUserAttempts();
        setAttempts(data.attempts || []);
      } catch (error) {
        console.error('Failed to load attempts:', error);
        setAttempts([]);
      } finally {
        setLoading(false);
      }
    };
    loadAttempts();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-white">Loading history...</div>
      </div>
    );
  }

  const totalQuizzes = attempts.length;
  const averageScore = attempts.length
    ? Math.round(attempts.reduce((sum, a) => sum + a.percentage, 0) / attempts.length)
    : 0;
  const recentActivity = "Today";

  return (
    <div className="min-h-screen bg-gray-900">
      <Header />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-2xl font-bold text-[#625FFF] mb-6">Quiz History</h1>

        {/* Top Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <div className="flex items-center space-x-4">
              <div className="w-12 h-12 rounded-lg bg-[#625FFF]/20 flex items-center justify-center">
          <svg className="w-6 h-6 text-[#625FFF]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
              </div>
              <div>
          <p className="text-gray-400 text-sm">Total Quizzes</p>
          <p className="text-3xl font-bold text-white">{totalQuizzes}</p>
              </div>
            </div>
          </div>

          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <div className="flex items-center space-x-4">
          <div className="w-12 h-12 rounded-lg bg-[#625FFF]/20 flex items-center justify-center">
          <svg className="w-6 h-6 text-[#625FFF]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          </div>
              <div>
          <p className="text-gray-400 text-sm">Average Score</p>
          <p className="text-3xl font-bold text-white">{averageScore}%</p>
              </div>
            </div>
          </div>

          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
             <div className="flex items-center space-x-4">
              <div className="w-12 h-12 rounded-lg bg-[#625FFF]/20 flex items-center justify-center">
          <svg className="w-6 h-6 text-[#625FFF]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
              </div>
               <div>
           <p className="text-gray-400 text-sm">Recent Activity</p>
           <p className="text-3xl font-bold text-white">{recentActivity}</p>
               </div>
             </div>
           </div>
        </div>

        {/* Filter & Sort Bar */}
        <div className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 mb-6 flex items-center justify-between">
          <div className="flex items-center space-x-2 text-gray-300">
            <span className="text-sm">Filter & Sort:</span>
            <button className="bg-gray-700 hover:bg-gray-600 text-white px-3 py-1 rounded-md text-sm">All Categories</button>
          </div>
          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-300">Sort by:</span>
            <button className="bg-gray-700 hover:bg-gray-600 text-white px-3 py-1 rounded-md text-sm">Most Recent</button>
          </div>
        </div>

        {/* History List */}
        <div className="space-y-4">
          {attempts.map((attempt) => (
            <div key={attempt.id} className="bg-gray-800 border border-gray-700 rounded-lg p-4 flex items-center justify-between">
              {/* Left: icon + title + meta */}
              <div className="flex items-center space-x-4">
                <div className="w-10 h-10 rounded-lg bg-[#625FFF]/20 flex items-center justify-center">
                  <svg className="w-6 h-6 text-[#625FFF]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div>
                  <h3 className="text-white font-semibold text-sm md:text-base">{attempt.quiz_title}</h3>
                  <div className="flex items-center space-x-4 text-sm text-gray-400 mt-1">
                    <div className="flex items-center space-x-1">
                      <span>{attempt.total_questions} questions</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <span>{formatDate(attempt.created_at)}</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Right: score + button */}
              <div className="flex items-center space-x-6">
                <div className="text-right">
                  <div className="text-white text-xl font-bold">{Math.round(attempt.percentage)}%</div>
                  <div className="text-gray-400 text-xs">{attempt.score}/{attempt.total_questions}</div>
                </div>
                <Link
                  href={`/quiz/results/${attempt.id}`}
                  className="bg-[#625FFF] hover:bg-[#544FFE] text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
                >
                  View Details
                </Link>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}