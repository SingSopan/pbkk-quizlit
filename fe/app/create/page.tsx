"use client";

import React, { useState, useRef } from "react";
import Link from "next/link";
import { generateQuizFromFile, generateQuizFromContent } from '../lib/quiz-service';
import Header from "../../components/Header";
// Removed import as we're using direct implementation

export default function CreateQuiz() {
  const [dragActive, setDragActive] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [quizDetails, setQuizDetails] = useState({
    title: "",
    description: "",
    questionCount: 10
  });
  const [showQuizDetails, setShowQuizDetails] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showErrorModal, setShowErrorModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [generatedQuiz, setGeneratedQuiz] = useState<any>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    const files = e.dataTransfer.files;
    if (files && files[0]) {
      handleFile(files[0]);
    }
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files[0]) {
      handleFile(files[0]);
    }
  };

  const handleFile = (file: File) => {
    const allowedTypes = ['application/pdf', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'text/plain'];
    if (allowedTypes.includes(file.type)) {
      setUploadedFile(file);
      setShowQuizDetails(true);
    } else {
      alert('Please upload a PDF, DOCX, or TXT file');
    }
  };

  const generateQuiz = async () => {
    if (!quizDetails.title || !quizDetails.description) {
      alert('Please fill in all quiz details');
      return;
    }

    setIsGenerating(true);

    try {
      let quiz;
      
      if (uploadedFile) {
        // Use file upload API
        quiz = await generateQuizFromFile(uploadedFile, {
          title: quizDetails.title,
          description: quizDetails.description,
          difficulty: "medium",
          questionCount: quizDetails.questionCount,
        });
      } else {
        // Use text generation API with placeholder content
        quiz = await generateQuizFromContent(
          "Please generate a general knowledge quiz about " + quizDetails.title,
          {
            title: quizDetails.title,
            description: quizDetails.description,
            difficulty: "medium",
            questionCount: quizDetails.questionCount,
          }
        );
      }

      setGeneratedQuiz(quiz);
      setShowSuccessModal(true);
    } catch (error: any) {
      // Silently handle errors - show modal instead of console logging
      if (error?.message?.includes('already exists') || error?.response?.status === 409) {
        setErrorMessage(`Quiz "${quizDetails.title}" already exists. Please choose a different name.`);
      } else {
        setErrorMessage('Failed to generate quiz. Please try again.');
      }
      setShowErrorModal(true);
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900">
      {/* Header */}
      <Header />

      {/* Success Modal */}
      {showSuccessModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4 border border-gray-700">
            <div className="flex flex-col items-center">
              {/* Success Icon */}
              <div className="w-20 h-20 mb-6 bg-green-500/20 rounded-full flex items-center justify-center">
                <svg className="w-12 h-12 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              
              {/* Text */}
              <h3 className="text-2xl font-bold text-white mb-2">
                Quiz Created Successfully!
              </h3>
              <p className="text-gray-400 text-center mb-6">
                Your quiz has been generated with {generatedQuiz?.questions?.length || 0} questions.
              </p>
              
              {/* Buttons */}
              <div className="flex gap-3 w-full">
                <button
                  onClick={() => window.location.href = '/dashboard'}
                  className="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-3 px-6 rounded-lg font-medium transition-colors"
                >
                  View Dashboard
                </button>
                <button
                  onClick={() => {
                    setShowSuccessModal(false);
                    setUploadedFile(null);
                    setQuizDetails({ title: '', description: '', questionCount: 10 });
                    setShowQuizDetails(false);
                  }}
                  className="flex-1 bg-gray-700 hover:bg-gray-600 text-white py-3 px-6 rounded-lg font-medium transition-colors"
                >
                  Create Another
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Error Modal */}
      {showErrorModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4 border border-red-900/50">
            <div className="flex flex-col items-center">
              {/* Error Icon */}
              <div className="w-20 h-20 mb-6 bg-red-500/20 rounded-full flex items-center justify-center">
                <svg className="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              
              {/* Text */}
              <h3 className="text-2xl font-bold text-white mb-2">
                Oops! Something Went Wrong
              </h3>
              <p className="text-gray-300 text-center mb-6">
                {errorMessage}
              </p>
              
              {/* Button */}
              <button
                onClick={() => setShowErrorModal(false)}
                className="w-full bg-red-600 hover:bg-red-700 text-white py-3 px-6 rounded-lg font-medium transition-colors"
              >
                Got it
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Loading Overlay */}
      {isGenerating && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4 border border-gray-700">
            <div className="flex flex-col items-center">
              {/* Spinner */}
              <div className="relative w-20 h-20 mb-6">
                <div className="absolute inset-0 border-4 border-blue-200 rounded-full"></div>
                <div className="absolute inset-0 border-4 border-transparent border-t-blue-600 rounded-full animate-spin"></div>
              </div>
              
              {/* Text */}
              <h3 className="text-xl font-semibold text-white mb-2">
                Generating Your Quiz
              </h3>
              <p className="text-gray-400 text-center">
                AI is analyzing your content and creating questions. This may take a moment...
              </p>
              
              {/* Progress dots */}
              <div className="flex gap-2 mt-6">
                <div className="w-2 h-2 bg-blue-600 rounded-full animate-bounce" style={{ animationDelay: '0ms' }}></div>
                <div className="w-2 h-2 bg-blue-600 rounded-full animate-bounce" style={{ animationDelay: '150ms' }}></div>
                <div className="w-2 h-2 bg-blue-600 rounded-full animate-bounce" style={{ animationDelay: '300ms' }}></div>
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Title */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-4">
            Upload <span className="text-blue-400">Your File</span>
          </h1>
          <p className="text-gray-300">
            Turn your PDFs, Docs, or Notes into Interactive quizzes instantly
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* File Upload Card */}
          <div className="bg-gray-800 rounded-lg p-8 border-2 border-gray-700 ">
            <div
              className={`relative border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                dragActive
                  ? 'border-blue-400 bg-blue-900/20'
                  : 'border-gray-600 bg-gray-700/50'
              }`}
              onDragEnter={handleDrag}
              onDragLeave={handleDrag}
              onDragOver={handleDrag}
              onDrop={handleDrop}
            >
              <input
                ref={fileInputRef}
                type="file"
                className="hidden"
                accept=".pdf,.docx,.txt"
                onChange={handleFileInput}
              />
              
              <div className="mb-4">
                <svg className="w-16 h-16 text-gray-400 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
              </div>
              
              <h3 className="text-xl font-semibold text-white mb-2">
                Drag & drop your file here
              </h3>
              <p className="text-gray-400 mb-4">
                or click to browse your files
              </p>
              
              <button
                onClick={() => fileInputRef.current?.click()}
                className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-lg transition-colors font-medium"
              >
                Choose File
              </button>
            </div>

            {/* Supported Formats */}
            <div className="mt-6">
              <p className="text-gray-400 text-sm text-center mb-3">Supported format:</p>
              <div className="flex justify-center space-x-4">
                <div className="flex items-center space-x-2 bg-gray-700 px-4 py-2 rounded-lg">
                  <span className="text-red-400 text-lg">📄</span>
                  <span className="text-gray-300 text-sm font-medium">PDF</span>
                </div>
              </div>
            </div>
          </div>

          {/* Study Materials Card */}
          <div className="bg-gray-800 rounded-lg p-8 border-2 border-gray-700 flex items-center justify-center">
            <div className="text-center">
              <div className="w-32 h-32 bg-blue-600/20 rounded-full mx-auto mb-6 flex items-center justify-center border-4 border-blue-600/30">
                <svg className="w-16 h-16 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">
                Upload Your Study Materials
              </h3>
              <p className="text-gray-400 leading-relaxed">
                Upload your study materials and we'll create engaging quizzes automatically using AI
              </p>
              
              {/* Features */}
              <div className="mt-6 space-y-2 text-sm text-gray-400">
                <div className="flex items-center justify-center gap-2">
                  <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  <span>AI-Powered Question Generation</span>
                </div>
                <div className="flex items-center justify-center gap-2">
                  <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  <span>Multiple Choice Format</span>
                </div>
                <div className="flex items-center justify-center gap-2">
                  <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  <span>Instant Quiz Creation</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Quiz Details Form */}
        {showQuizDetails && (
          <div className="mt-8 bg-gray-800 rounded-lg p-6">
            <h2 className="text-2xl font-bold text-white mb-6">Quiz Details</h2>
            <p className="text-gray-400 mb-6">Give your quiz a name and description</p>
            
            {/* Uploaded File Display */}
            {uploadedFile && (
              <div className="mb-6 bg-gray-700 rounded-lg p-4 flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
                    <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-white font-medium">{uploadedFile.name}</p>
                    <p className="text-gray-400 text-sm">{(uploadedFile.size / 1024).toFixed(2)} KB</p>
                  </div>
                </div>
                <button
                  onClick={() => {
                    setUploadedFile(null);
                    setShowQuizDetails(false);
                  }}
                  className="text-red-400 hover:text-red-300 p-2"
                  title="Remove file"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            )}
            
            <div className="grid grid-cols-1 gap-6">
              <div>
                <label className="block text-gray-300 text-sm font-medium mb-2">
                  Quiz Title
                </label>
                <input
                  type="text"
                  value={quizDetails.title}
                  onChange={(e) => setQuizDetails({...quizDetails, title: e.target.value})}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                  placeholder="Enter quiz title"
                />
              </div>
              
              <div>
                <label className="block text-gray-300 text-sm font-medium mb-2">
                  Quiz Description
                </label>
                <textarea
                  value={quizDetails.description}
                  onChange={(e) => setQuizDetails({...quizDetails, description: e.target.value})}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                  rows={3}
                  placeholder="Enter quiz description"
                />
              </div>
              
              <div>
                <label className="block text-gray-300 text-sm font-medium mb-2">
                  Number of Questions
                </label>
                <input
                  type="number"
                  min="5"
                  max="15"
                  value={quizDetails.questionCount}
                  onChange={(e) => setQuizDetails({...quizDetails, questionCount: parseInt(e.target.value) || 10})}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                  placeholder="Enter number of questions (5-15)"
                />
                <p className="text-gray-400 text-sm mt-1">Choose between 5 and 15 questions</p>
              </div>
            </div>
            
            <button
              onClick={generateQuiz}
              disabled={!quizDetails.title || !quizDetails.description || isGenerating}
              className="mt-6 w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white py-3 px-6 rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
            >
              {isGenerating ? (
                <>
                  <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Generating Quiz...
                </>
              ) : 'Generate Quiz'}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}