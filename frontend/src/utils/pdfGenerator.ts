import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'
import type { AnalysisResult } from '../types/analysis'

export const generateAnalysisPDF = (result: AnalysisResult, fileName: string) => {
    const doc = new jsPDF()
    const primaryColor: [number, number, number] = [29, 78, 216] // Indigo-700
    const secondaryColor: [number, number, number] = [71, 85, 105] // Slate-600

    const originalName = fileName.replace(/\.[^/.]+$/, "")
    const pdfFileName = `${originalName}_analysis.pdf`

    // --- Header ---
    doc.setFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
    doc.rect(0, 0, 210, 40, 'F')

    doc.setTextColor(255, 255, 255)
    doc.setFontSize(22)
    doc.setFont('helvetica', 'bold')
    doc.text('CV EVALUATION REPORT', 15, 25)

    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    doc.text(`Generated on: ${new Date().toLocaleDateString()}`, 15, 33)

    // --- Overall Score Circle/Badge ---
    doc.setFillColor(255, 255, 255)
    doc.roundedRect(160, 10, 35, 20, 3, 3, 'F')
    doc.setTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
    doc.setFontSize(16)
    doc.text(`${(result.matching_score.overall * 10).toFixed(0)}%`, 177.5, 23, { align: 'center' })
    doc.setFontSize(8)
    doc.text('MATCH SCORE', 177.5, 28, { align: 'center' })

    let currentY = 50

    // --- Candidate Highlights Section ---
    doc.setTextColor(0, 0, 0)
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Candidate Highlights', 15, currentY)
    currentY += 8

    autoTable(doc, {
        startY: currentY,
        head: [['Strategic Metric', 'Value']],
        body: [
            ['File Name', fileName],
            ['Evaluation ID', result.processing_metadata.request_id || 'N/A'],
        ],
        theme: 'striped',
        headStyles: { fillColor: primaryColor },
        margin: { left: 15, right: 15 }
    })
    currentY = (doc as any).lastAutoTable.finalY + 10

    // --- Executive Summary ---
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Executive Summary', 15, currentY)
    currentY += 6

    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    const summaryLines = doc.splitTextToSize(result.executive_summary, 180)
    doc.text(summaryLines, 15, currentY)
    currentY += (summaryLines.length * 5) + 10

    // --- Technical Skills Table ---
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Key Technical Alignment', 15, currentY)
    currentY += 6

    const skillData = result.technical_skills.map(s => [s.name, `${(s.score * 10).toFixed(0)}%`])
    autoTable(doc, {
        startY: currentY,
        head: [['Technology / Competency', 'Proficiency Score']],
        body: skillData,
        theme: 'grid',
        headStyles: { fillColor: secondaryColor },
        margin: { left: 15, right: 15 }
    })
    currentY = (doc as any).lastAutoTable.finalY + 10

    // Check for page break
    if (currentY > 230) {
        doc.addPage()
        currentY = 20
    }

    // --- Dimension Breakdown ---
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Matching Dimensions', 15, currentY)
    currentY += 6

    const dimensionData = [
        ['Skills Alignment', `${(result.matching_score.skills_alignment.score * 10).toFixed(0)}%`, result.matching_score.skills_alignment.evidence[0] || ''],
        ['Exp Relevance', `${(result.matching_score.experience_relevance.score * 10).toFixed(0)}%`, result.matching_score.experience_relevance.evidence[0] || ''],
        ['Seniority Fit', `${(result.matching_score.seniority_fit.score * 10).toFixed(0)}%`, result.matching_score.seniority_fit.evidence[0] || ''],
        ['Domain Knowledge', `${(result.matching_score.domain_knowledge.score * 10).toFixed(0)}%`, result.matching_score.domain_knowledge.evidence[0] || ''],
    ]

    autoTable(doc, {
        startY: currentY,
        head: [['Dimension', 'Score', 'Primary Reasoning']],
        body: dimensionData,
        theme: 'striped',
        styles: { fontSize: 8 },
        columnStyles: { 2: { cellWidth: 100 } },
        headStyles: { fillColor: primaryColor }
    })
    currentY = (doc as any).lastAutoTable.finalY + 10

    // --- Strengths & Gaps ---
    if (currentY > 230) {
        doc.addPage()
        currentY = 20
    }

    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Key Strengths', 15, currentY)
    currentY += 6
    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    result.strengths.slice(0, 3).forEach(s => {
        const lines = doc.splitTextToSize(`• ${s.title}: ${s.description}`, 180)
        doc.text(lines, 15, currentY)
        currentY += (lines.length * 5) + 2
    })
    currentY += 5

    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Critical Gaps', 15, currentY)
    currentY += 6
    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    result.gaps.slice(0, 3).forEach(g => {
        const lines = doc.splitTextToSize(`• [${g.severity}] ${g.title}: ${g.description}`, 180)
        doc.text(lines, 15, currentY)
        currentY += (lines.length * 5) + 2
    })
    currentY += 10

    // --- Footer ---
    const pageCount = (doc as any).internal.getNumberOfPages()
    for (let i = 1; i <= pageCount; i++) {
        doc.setPage(i)
        doc.setFontSize(8)
        doc.setTextColor(150)
        doc.text(`Page ${i} of ${pageCount} | Confidential Recruitment Intelligence`, 105, 285, { align: 'center' })
    }

    doc.save(pdfFileName)
}
