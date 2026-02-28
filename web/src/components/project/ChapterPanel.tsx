import { useState } from 'react';
import { Button, Card, List, Tag, message, Modal, InputNumber, Spin } from 'antd';
import { PlayCircleOutlined, CheckOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { chapterApi } from '../../api';
import type { Project, Chapter } from '../../types';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

interface ChapterPanelProps {
  project: Project;
}

function ChapterPanel({ project }: ChapterPanelProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [selectedChapter, setSelectedChapter] = useState<Chapter | null>(null);
  const [detailModalOpen, setDetailModalOpen] = useState(false);

  // 获取章节列表
  const { data: chaptersData, isLoading } = useQuery({
    queryKey: ['chapters', project.id, page],
    queryFn: () =>
      chapterApi.list(project.id, { page, page_size: 20 }).then((res) => res.data.data),
    enabled: !!project.id,
  });

  const chapters = chaptersData?.chapters || [];
  const total = chaptersData?.total || 0;

  // 生成章节内容
  const generateMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.generateContent(project.id, chapterNumber, { overwrite: false }),
    onSuccess: () => {
      message.success('章节内容生成成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
    },
    onError: (err) => {
      if (err.message?.includes('已有内容')) {
        message.error('章节已有内容，请先清空或选择覆盖');
      } else {
        message.error('生成失败');
      }
    },
  });

  const handleGenerate = (chapterNumber: number) => {
    Modal.confirm({
      title: '生成章节内容',
      content: `确定要生成第 ${chapterNumber} 章的内容吗？`,
      onOk: () => {
        generateMutation.mutate(chapterNumber);
      },
    });
  };

  const handleChapterClick = (chapter: Chapter) => {
    setSelectedChapter(chapter);
    setDetailModalOpen(true);
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      not_started: { color: 'default', text: '未开始' },
      draft: { color: 'processing', text: '草稿' },
      completed: { color: 'success', text: '已完成' },
    };
    const { color, text } = statusMap[status] || { color: 'default', text: status };
    return <Tag color={color}>{text}</Tag>;
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <div>
          <h3 className="text-lg font-semibold">章节写作</h3>
          <p className="text-gray-500 text-sm">
            进度: {project.completed_chapters}/{project.chapter_count} 章 |
            总字数: {project.total_words.toLocaleString()} 字
          </p>
        </div>
      </div>

      {isLoading ? (
        <div className="text-center py-12">
          <Spin />
        </div>
      ) : chapters.length === 0 ? (
        <Card className="text-center py-8">
          <p className="mb-4">暂无章节，请先生成章节大纲</p>
          <Button type="primary" onClick={() => navigate(`/projects/${project.id}`)}>
            前往生成大纲
          </Button>
        </Card>
      ) : (
        <List
          grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
          dataSource={chapters}
          pagination={{
            total,
            pageSize: 20,
            current: page,
            onChange: setPage,
          }}
          renderItem={(chapter) => (
            <List.Item>
              <Card
                size="small"
                hoverable
                onClick={() => handleChapterClick(chapter)}
                className={chapter.is_finalized ? 'border-green-500' : ''}
              >
                <div className="space-y-2">
                  <div className="flex justify-between items-start">
                    <span className="font-semibold">第 {chapter.chapter_number} 章</span>
                    {getStatusTag(chapter.status)}
                  </div>
                  {chapter.title && (
                    <div className="text-sm font-medium truncate">{chapter.title}</div>
                  )}
                  {chapter.word_count > 0 && (
                    <div className="text-xs text-gray-500">
                      {chapter.word_count} 字
                    </div>
                  )}
                  <div className="flex gap-1">
                    {chapter.status === 'not_started' && (
                      <Button
                        size="small"
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        onClick={(e) => {
                          e.stopPropagation();
                          handleGenerate(chapter.chapter_number);
                        }}
                        loading={generateMutation.isPending}
                      >
                        生成
                      </Button>
                    )}
                    {chapter.is_finalized && (
                      <Tag icon={<CheckOutlined />} color="success">
                        已定稿
                      </Tag>
                    )}
                  </div>
                </div>
              </Card>
            </List.Item>
          )}
        />
      )}

      {/* 章节详情弹窗 */}
      <Modal
        title={`第 ${selectedChapter?.chapter_number} 章${selectedChapter?.title ? ' - ' + selectedChapter.title : ''}`}
        open={detailModalOpen}
        onCancel={() => {
          setDetailModalOpen(false);
          setSelectedChapter(null);
        }}
        footer={null}
        width={800}
      >
        {selectedChapter && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">章节号</label>
              <InputNumber value={selectedChapter.chapter_number} disabled className="w-full" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">章节标题</label>
              <input
                type="text"
                value={selectedChapter.title || ''}
                className="w-full px-3 py-2 border rounded"
                placeholder="章节标题"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">章节内容</label>
              <textarea
                value={selectedChapter.content || ''}
                rows={15}
                className="w-full px-3 py-2 border rounded"
                placeholder="章节内容..."
              />
            </div>
            <div className="text-sm text-gray-500">
              字数: {selectedChapter.word_count}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}

export default ChapterPanel;
