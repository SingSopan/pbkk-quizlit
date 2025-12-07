"use client";

import { useState } from "react";
import Header from "@/components/Header";

export default function Documentation() {
  const [activeSection, setActiveSection] = useState<string>("intro");

  const sections = [
    { id: "intro", title: "Pengenalan" },
    { id: "setup", title: "Setup Awal" },
    { id: "create-quiz", title: "Membuat Quiz" },
    { id: "take-quiz", title: "Mengerjakan Quiz" },
    { id: "view-results", title: "Melihat Hasil" },
    { id: "history", title: "Riwayat Quiz" },
  ];

  return (
    <div className="min-h-screen bg-slate-900">
      <Header />
      
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <h1 className="text-4xl font-bold text-slate-100 mb-8 text-center">
            📚 Dokumentasi Quizlit
          </h1>

          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Sidebar Navigation */}
            <div className="lg:col-span-1">
              <div className="bg-slate-800/80 rounded-lg shadow-md border border-slate-700/50 p-4 sticky top-4">
                <h2 className="text-lg font-semibold mb-4 text-slate-200">
                  Daftar Isi
                </h2>
                <nav className="space-y-2">
                  {sections.map((section) => (
                    <button
                      key={section.id}
                      onClick={() => setActiveSection(section.id)}
                      className={`w-full text-left px-4 py-2 rounded-lg transition-colors ${
                        activeSection === section.id
                          ? "bg-indigo-600 text-white"
                          : "text-slate-300 hover:bg-slate-700/50"
                      }`}
                    >
                      {section.title}
                    </button>
                  ))}
                </nav>
              </div>
            </div>

            {/* Main Content */}
            <div className="lg:col-span-3">
              <div className="bg-slate-800/80 rounded-lg shadow-md border border-slate-700/50 p-8">
                {activeSection === "intro" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      🎯 Pengenalan Quizlit
                    </h2>
                    
                    {/* Video Tutorial */}
                    <div className="bg-slate-700/50 rounded-lg overflow-hidden">
                      <div className="aspect-video">
                        <iframe
                          src="https://drive.google.com/file/d/1j4B4YR7jEKuik5ZglNCh1WA_rYSVTGLu/preview"
                          className="w-full h-full"
                          allow="autoplay"
                          allowFullScreen
                        ></iframe>
                      </div>
                      <div className="p-4 bg-slate-800">
                        <p className="text-sm text-slate-300 flex items-center gap-2">
                          <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span className="font-medium">🎥 Video Tutorial: Pengenalan Aplikasi Quizlit</span>
                        </p>
                      </div>
                    </div>
                    
                    <p className="text-slate-300 leading-relaxed">
                      Quizlit adalah aplikasi quiz interaktif yang memungkinkan Anda membuat dan mengerjakan quiz dengan mudah. 
                      Aplikasi ini mendukung berbagai jenis pertanyaan dan dapat menghasilkan quiz otomatis dari dokumen PDF.
                    </p>
                    
                    <div className="bg-blue-500/10 border-l-4 border-blue-500 p-4 rounded">
                      <h3 className="font-semibold text-blue-400 mb-2">✨ Fitur Utama:</h3>
                      <ul className="list-disc list-inside text-slate-300 space-y-1">
                        <li>Upload PDF dan generate quiz otomatis dengan AI</li>
                        <li>Sistem penilaian otomatis</li>
                        <li>Riwayat quiz dan hasil detail</li>
                        <li>Dashboard pengguna yang informatif</li>
                      </ul>
                    </div>
                  </div>
                )}

                {activeSection === "setup" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      ⚙️ Setup Awal
                    </h2>
                    
                    <div className="space-y-4">
                      <div>
                        <h3 className="text-xl font-semibold text-slate-200 mb-3">
                          1. Registrasi Akun
                        </h3>
                        <div className="bg-slate-700/30 p-4 rounded-lg">
                          <ol className="list-decimal list-inside space-y-2 text-slate-300">
                            <li>Klik menu <strong>Register</strong> di navigation bar</li>
                            <li>Isi form dengan data Anda:
                              <ul className="list-disc list-inside ml-6 mt-2 space-y-1">
                                <li>Nama lengkap</li>
                                <li>Email aktif</li>
                                <li>Password (minimal 6 karakter)</li>
                              </ul>
                            </li>
                            <li>Klik tombol <strong>Sign Up</strong></li>
                            <li>Anda akan otomatis diarahkan ke halaman login</li>
                          </ol>
                        </div>
                      </div>

                      <div>
                        <h3 className="text-xl font-semibold text-slate-200 mb-3">
                          2. Login ke Aplikasi
                        </h3>
                        <div className="bg-slate-700/30 p-4 rounded-lg">
                          <ol className="list-decimal list-inside space-y-2 text-slate-300">
                            <li>Klik menu <strong>Login</strong> di navigation bar</li>
                            <li>Masukkan email dan password Anda</li>
                            <li>Klik tombol <strong>Sign In</strong></li>
                            <li>Setelah berhasil login, Anda akan diarahkan ke Dashboard</li>
                          </ol>
                        </div>
                      </div>
                    </div>

                    <div className="bg-yellow-500/10 border-l-4 border-yellow-500 p-4 rounded">
                      <p className="text-yellow-400">
                        <strong>💡 Tips:</strong> Pastikan Anda menggunakan email yang valid karena akan digunakan untuk autentikasi.
                      </p>
                    </div>
                  </div>
                )}

                {activeSection === "create-quiz" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      📝 Membuat Quiz
                    </h2>
                    
                    {/* Video Tutorial */}
                    <div className="bg-gradient-to-br from-purple-50 to-pink-50 rounded-lg overflow-hidden">
                      <div className="aspect-video">
                        <iframe
                          src="https://drive.google.com/file/d/1j4B4YR7jEKuik5ZglNCh1WA_rYSVTGLu/preview"
                          className="w-full h-full"
                          allow="autoplay"
                          allowFullScreen
                        ></iframe>
                      </div>
                      <div className="p-4 bg-slate-800">
                        <p className="text-sm text-purple-400 flex items-center gap-2">
                          <svg className="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span className="font-medium">🎥 Video Tutorial: Cara Membuat Quiz dari PDF</span>
                        </p>
                      </div>
                    </div>
                    
                    <div className="space-y-4">
                      <div>
                        <h3 className="text-xl font-semibold text-slate-200 mb-3">
                          Langkah-langkah:
                        </h3>
                        <div className="space-y-4">
                          <div className="bg-purple-500/10 p-4 rounded-lg border-l-4 border-purple-500">
                            <h4 className="font-semibold text-purple-400 mb-2">
                              1️⃣ Akses Halaman Create Quiz
                            </h4>
                            <p className="text-slate-300">
                              Klik tombol <strong>"Create New Quiz"</strong> di Dashboard atau menu <strong>"Create Quiz"</strong> di navigation bar.
                            </p>
                          </div>

                          <div className="bg-blue-500/10 p-4 rounded-lg border-l-4 border-blue-500">
                            <h4 className="font-semibold text-blue-400 mb-2">
                              2️⃣ Upload File PDF
                            </h4>
                            <ul className="list-disc list-inside text-slate-300 space-y-1">
                              <li>Klik area upload atau drag & drop file PDF Anda</li>
                              <li>File maksimal 20MB</li>
                              <li>Pastikan PDF berisi materi yang akan dijadikan quiz</li>
                              <li>Tunggu hingga proses upload selesai</li>
                            </ul>
                          </div>

                          <div className="bg-green-500/10 p-4 rounded-lg border-l-4 border-green-500">
                            <h4 className="font-semibold text-green-400 mb-2">
                              3️⃣ Konfigurasi Quiz
                            </h4>
                            <p className="text-slate-300 mb-2">
                              Setelah upload berhasil, isi detail quiz:
                            </p>
                            <ul className="list-disc list-inside text-slate-300 space-y-1">
                              <li><strong>Judul Quiz:</strong> Berikan nama yang jelas</li>
                              <li><strong>Deskripsi:</strong> Jelaskan topik atau tujuan quiz (opsional)</li>
                              <li><strong>Jumlah Soal:</strong> Pilih berapa banyak soal yang ingin dibuat (5-15 soal)</li>
                            </ul>
                          </div>

                          <div className="bg-orange-500/10 p-4 rounded-lg border-l-4 border-orange-500">
                            <h4 className="font-semibold text-orange-400 mb-2">
                              4️⃣ Generate Quiz dengan AI
                            </h4>
                            <p className="text-slate-300 mb-2">
                              Klik tombol <strong>"Generate Quiz"</strong>:
                            </p>
                            <ul className="list-disc list-inside text-slate-300 space-y-1">
                              <li>AI akan menganalisis konten PDF Anda</li>
                              <li>Proses membutuhkan waktu 10-30 detik</li>
                              <li>Soal akan dibuat otomatis berdasarkan materi PDF</li>
                            </ul>
                          </div>

                          <div className="bg-pink-500/10 p-4 rounded-lg border-l-4 border-pink-500">
                            <h4 className="font-semibold text-pink-400 mb-2">
                              5️⃣ Selesai! ✅
                            </h4>
                            <p className="text-slate-300">
                              Setelah quiz berhasil dibuat, Anda akan diarahkan ke halaman detail quiz. 
                              Quiz sudah siap dikerjakan atau dibagikan!
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className="bg-red-500/10 border-l-4 border-red-500 p-4 rounded">
                      <p className="text-red-400">
                        <strong>⚠️ Perhatian:</strong> Pastikan koneksi internet stabil saat generate quiz karena proses ini menggunakan AI yang memerlukan koneksi ke server.
                      </p>
                    </div>
                  </div>
                )}

                {activeSection === "take-quiz" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      ✍️ Mengerjakan Quiz
                    </h2>
                    
                    {/* Video Tutorial */}
                    <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-lg overflow-hidden">
                      <div className="aspect-video">
                        <iframe
                          src="https://drive.google.com/file/d/1j4B4YR7jEKuik5ZglNCh1WA_rYSVTGLu/preview"
                          className="w-full h-full"
                          allow="autoplay"
                          allowFullScreen
                        ></iframe>
                      </div>
                      <div className="p-4 bg-slate-800">
                        <p className="text-sm text-blue-400 flex items-center gap-2">
                          <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span className="font-medium">🎥 Video Tutorial: Cara Mengerjakan Quiz</span>
                        </p>
                      </div>
                    </div>
                    
                    <div className="space-y-4">
                      <div className="bg-indigo-500/10 p-4 rounded-lg border-l-4 border-indigo-500">
                        <h4 className="font-semibold text-indigo-400 mb-2">
                          1️⃣ Mulai Quiz
                        </h4>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li>Dari Dashboard, pilih quiz yang ingin dikerjakan</li>
                          <li>Klik tombol <strong>"Take Quiz"</strong> atau <strong>"Start Quiz"</strong></li>
                          <li>Baca instruksi dengan teliti sebelum memulai</li>
                        </ul>
                      </div>

                      <div className="bg-cyan-500/10 p-4 rounded-lg border-l-4 border-cyan-500">
                        <h4 className="font-semibold text-cyan-400 mb-2">
                          2️⃣ Menjawab Soal
                        </h4>
                        <div className="text-slate-300 space-y-3">
                          <div>
                            <p className="font-medium mb-1">📌 Multiple Choice:</p>
                            <ul className="list-disc list-inside ml-4 space-y-1">
                              <li>Pilih salah satu opsi jawaban yang tersedia</li>
                              <li>Klik pada pilihan untuk menyeleksi</li>
                            </ul>
                          </div>
                        </div>
                      </div>

                      <div className="bg-teal-500/10 p-4 rounded-lg border-l-4 border-teal-500">
                        <h4 className="font-semibold text-teal-400 mb-2">
                          3️⃣ Navigasi Soal
                        </h4>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li>Gunakan tombol <strong>"Next"</strong> untuk ke soal berikutnya</li>
                          <li>Gunakan tombol <strong>"Previous"</strong> untuk kembali ke soal sebelumnya</li>
                          <li>Indikator progress menunjukkan posisi soal saat ini</li>
                          <li>Anda bisa kembali ke soal sebelumnya untuk mengubah jawaban</li>
                        </ul>
                      </div>

                      <div className="bg-emerald-500/10 p-4 rounded-lg border-l-4 border-emerald-500">
                        <h4 className="font-semibold text-emerald-400 mb-2">
                          4️⃣ Submit Quiz
                        </h4>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li>Setelah semua soal dijawab, klik tombol <strong>"Submit Quiz"</strong></li>
                          <li>Akan muncul konfirmasi sebelum submit</li>
                          <li>Pastikan semua jawaban sudah terisi</li>
                          <li>Setelah submit, jawaban tidak bisa diubah lagi</li>
                        </ul>
                      </div>

                      <div className="bg-lime-500/10 p-4 rounded-lg border-l-4 border-lime-500">
                        <h4 className="font-semibold text-lime-400 mb-2">
                          5️⃣ Lihat Hasil
                        </h4>
                        <p className="text-slate-300">
                          Setelah submit, Anda akan langsung diarahkan ke halaman hasil yang menampilkan:
                        </p>
                        <ul className="list-disc list-inside text-slate-300 space-y-1 mt-2">
                          <li>Skor total Anda</li>
                          <li>Jumlah jawaban benar dan salah</li>
                          <li>Pembahasan untuk setiap soal</li>
                          <li>Jawaban yang benar</li>
                        </ul>
                      </div>
                    </div>

                    <div className="bg-blue-500/10 border-l-4 border-blue-500 p-4 rounded">
                      <p className="text-blue-400">
                        <strong>💡 Tips:</strong> Baca soal dengan cermat sebelum menjawab. Untuk soal essay, usahakan jawaban yang komprehensif dan terstruktur.
                      </p>
                    </div>
                  </div>
                )}

                {activeSection === "view-results" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      📊 Melihat Hasil Quiz
                    </h2>
                    
                    {/* Video Tutorial */}
                    <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg overflow-hidden">
                      <div className="aspect-video">
                        <iframe
                          src="https://drive.google.com/file/d/1j4B4YR7jEKuik5ZglNCh1WA_rYSVTGLu/preview"
                          className="w-full h-full"
                          allow="autoplay"
                          allowFullScreen
                        ></iframe>
                      </div>
                      <div className="p-4 bg-slate-800">
                        <p className="text-sm text-green-400 flex items-center gap-2">
                          <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span className="font-medium">🎥 Video Tutorial: Memahami Hasil Quiz</span>
                        </p>
                      </div>
                    </div>
                    
                    <div className="space-y-4">
                      <div className="bg-violet-500/10 p-4 rounded-lg border-l-4 border-violet-500">
                        <h4 className="font-semibold text-violet-400 mb-2">
                          Halaman Hasil Quiz
                        </h4>
                        <p className="text-slate-300 mb-3">
                          Setelah menyelesaikan quiz, Anda akan melihat halaman hasil yang komprehensif:
                        </p>
                        <div className="space-y-2 text-slate-300">
                          <div className="bg-slate-700/30 p-3 rounded">
                            <p className="font-medium">🎯 Skor & Statistik</p>
                            <ul className="list-disc list-inside ml-4 text-sm space-y-1 mt-1">
                              <li>Skor total dalam persentase</li>
                              <li>Jumlah soal yang dijawab benar</li>
                              <li>Jumlah soal yang dijawab salah</li>
                              <li>Total soal dalam quiz</li>
                            </ul>
                          </div>
                          
                          <div className="bg-slate-700/30 p-3 rounded">
                            <p className="font-medium">📝 Detail Jawaban</p>
                            <ul className="list-disc list-inside ml-4 text-sm space-y-1 mt-1">
                              <li>Setiap soal ditampilkan dengan lengkap</li>
                              <li>Jawaban Anda ditandai dengan jelas</li>
                              <li>Jawaban yang benar ditampilkan</li>
                              <li>Indikator benar/salah dengan warna</li>
                            </ul>
                          </div>
                        </div>
                      </div>

                      <div className="bg-fuchsia-500/10 p-4 rounded-lg border-l-4 border-fuchsia-500">
                        <h4 className="font-semibold text-fuchsia-400 mb-2">
                          Aksi yang Tersedia
                        </h4>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li><strong>Back to Dashboard:</strong> Kembali ke halaman utama</li>
                          <li><strong>Retake Quiz:</strong> Mengerjakan quiz yang sama lagi</li>
                          <li><strong>View History:</strong> Melihat semua riwayat quiz Anda</li>
                        </ul>
                      </div>
                    </div>

                    <div className="bg-green-500/10 border-l-4 border-green-500 p-4 rounded">
                      <p className="text-green-400">
                        <strong>✨ Catatan:</strong> Hasil quiz disimpan secara permanen dan dapat diakses kapan saja melalui halaman History.
                      </p>
                    </div>
                  </div>
                )}

                {activeSection === "history" && (
                  <div className="space-y-6">
                    <h2 className="text-3xl font-bold text-slate-100 mb-4">
                      📜 Riwayat Quiz
                    </h2>
                    
                    <div className="space-y-4">
                      <div className="bg-rose-500/10 p-4 rounded-lg border-l-4 border-rose-500">
                        <h4 className="font-semibold text-rose-400 mb-2">
                          Mengakses Riwayat
                        </h4>
                        <ol className="list-decimal list-inside text-slate-300 space-y-2">
                          <li>Klik menu <strong>"History"</strong> di navigation bar</li>
                          <li>Atau klik <strong>"View Quiz History"</strong> di Dashboard</li>
                          <li>Anda akan melihat daftar semua quiz yang pernah dikerjakan</li>
                        </ol>
                      </div>

                      <div className="bg-amber-500/10 p-4 rounded-lg border-l-4 border-amber-500">
                        <h4 className="font-semibold text-amber-400 mb-2">
                          Informasi dalam Riwayat
                        </h4>
                        <p className="text-slate-300 mb-2">
                          Setiap entri riwayat menampilkan:
                        </p>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li><strong>Judul Quiz:</strong> Nama quiz yang dikerjakan</li>
                          <li><strong>Tanggal:</strong> Kapan quiz dikerjakan</li>
                          <li><strong>Skor:</strong> Nilai yang didapat</li>
                          <li><strong>Status:</strong> Completed atau In Progress</li>
                          <li><strong>Jumlah Soal:</strong> Total pertanyaan dalam quiz</li>
                        </ul>
                      </div>

                      <div className="bg-sky-500/10 p-4 rounded-lg border-l-4 border-sky-500">
                        <h4 className="font-semibold text-sky-400 mb-2">
                          Aksi pada Riwayat
                        </h4>
                        <ul className="list-disc list-inside text-slate-300 space-y-1">
                          <li><strong>View Details:</strong> Melihat detail hasil quiz</li>
                          <li><strong>Review Answers:</strong> Meninjau semua jawaban dan pembahasan</li>
                          <li><strong>Retake:</strong> Mengerjakan quiz yang sama lagi</li>
                        </ul>
                      </div>

                      <div className="bg-purple-500/10 p-4 rounded-lg border-l-4 border-purple-500">
                        <h4 className="font-semibold text-purple-400 mb-2">
                          Filter & Pencarian
                        </h4>
                        <p className="text-slate-300">
                          Gunakan fitur pencarian dan filter untuk menemukan riwayat quiz tertentu:
                        </p>
                        <ul className="list-disc list-inside text-slate-300 space-y-1 mt-2">
                          <li>Cari berdasarkan judul quiz</li>
                          <li>Filter berdasarkan tanggal</li>
                          <li>Urutkan berdasarkan skor atau tanggal</li>
                        </ul>
                      </div>
                    </div>

                    <div className="bg-purple-500/10 border-l-4 border-purple-500 p-4 rounded">
                      <p className="text-purple-400">
                        <strong>📌 Info:</strong> Riwayat quiz berguna untuk melacak progress belajar Anda dan mengidentifikasi area yang perlu diperbaiki.
                      </p>
                    </div>
                  </div>
                )}

                {/* Quick Tips Section - Always visible at bottom */}
                <div className="mt-8 pt-6 border-t border-slate-700">
                  <h3 className="text-xl font-semibold text-slate-200 mb-4">
                    🚀 Quick Tips
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="bg-purple-500/10 p-4 rounded-lg border border-purple-500/20">
                      <p className="text-sm text-slate-300">
                        <strong className="text-purple-400">💾 Auto Save:</strong> Jawaban Anda otomatis tersimpan saat navigasi antar soal
                      </p>
                    </div>
                    <div className="bg-blue-500/10 p-4 rounded-lg border border-blue-500/20">
                      <p className="text-sm text-slate-300">
                        <strong className="text-blue-400">🔄 Retake:</strong> Anda bisa mengulang quiz berapa kali pun untuk meningkatkan skor
                      </p>
                    </div>
                    <div className="bg-green-500/10 p-4 rounded-lg border border-green-500/20">
                      <p className="text-sm text-slate-300">
                        <strong className="text-green-400">📱 Responsive:</strong> Aplikasi dapat diakses dari desktop, tablet, atau smartphone
                      </p>
                    </div>
                    <div className="bg-orange-500/10 p-4 rounded-lg border border-orange-500/20">
                      <p className="text-sm text-slate-300">
                        <strong className="text-orange-400">🔒 Secure:</strong> Data Anda aman dan terlindungi dengan enkripsi
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
