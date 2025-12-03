import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <header className="bg-slate-800 border-b border-slate-700 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-5">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                QuizLit
              </h1>
            </div>
            <div className="flex items-center space-x-4">
              <Link
                href="/login"
                className="text-slate-300 hover:text-indigo-400 px-4 py-2 rounded-xl text-sm font-semibold transition-colors"
              >
                Login
              </Link>
              <Link
                href="/register"
                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2.5 rounded-xl text-sm font-semibold transition-all shadow-lg hover:shadow-xl"
              >
                Sign Up
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center">
          <h2 className="text-5xl font-extrabold text-white sm:text-6xl md:text-7xl leading-tight mb-6">
            Create Amazing Quizzes
          </h2>
          <p className="max-w-2xl mx-auto text-lg text-slate-400 sm:text-xl leading-relaxed mb-10">
            Generate interactive quizzes in seconds. Perfect for educators, trainers, and anyone who wants to test knowledge in a fun way.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/register"
              className="inline-flex items-center justify-center px-8 py-4 text-base font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 shadow-lg hover:shadow-xl transition-all"
            >
              Get Started Free
            </Link>
            <Link
              href="/login"
              className="inline-flex items-center justify-center px-8 py-4 text-base font-medium rounded-lg text-slate-200 bg-slate-800 border border-slate-600 hover:bg-slate-700 transition-all"
            >
              Sign In
            </Link>
          </div>
        </div>

        {/* Everything You Need Section */}
        <div className="mt-32">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-extrabold text-white sm:text-5xl mb-4">
              Everything you need
            </h2>
            <p className="max-w-3xl mx-auto text-lg text-slate-400">
              From document upload to detailed analytics, we've got every step covered
            </p>
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {/* Quick Generation */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-indigo-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Quick Generation</h3>
              <p className="text-slate-400 leading-relaxed">
                Create quizzes in minutes with our intuitive interface and smart question generation.
              </p>
            </div>

            {/* Easy to Use */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-emerald-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Easy to Use</h3>
              <p className="text-slate-400 leading-relaxed">
                Intuitive design makes creating and managing quizzes a breeze for everyone.
              </p>
            </div>

            {/* Customizable */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-purple-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Customizable</h3>
              <p className="text-slate-400 leading-relaxed">
                Customize your quizzes with different question types and difficulty levels.
              </p>
            </div>

            {/* Smart AI */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-amber-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">AI Powered</h3>
              <p className="text-slate-400 leading-relaxed">
                Leverage AI to generate relevant questions from your documents automatically.
              </p>
            </div>

            {/* Analytics */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-cyan-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Detailed Analytics</h3>
              <p className="text-slate-400 leading-relaxed">
                Track performance and get insights with comprehensive analytics dashboard.
              </p>
            </div>

            {/* Secure */}
            <div className="relative bg-slate-800/50 rounded-xl p-8 hover:bg-slate-800 transition-all border border-slate-700/50 hover:border-indigo-500/30">
              <div className="w-12 h-12 bg-violet-500/10 rounded-lg flex items-center justify-center mb-6">
                <svg className="h-6 w-6 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Secure & Private</h3>
              <p className="text-slate-400 leading-relaxed">
                Your data is encrypted and secure. We take privacy seriously.
              </p>
            </div>
          </div>
        </div>


      </main>
    </div>
  );
}
