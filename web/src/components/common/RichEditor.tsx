import { useEffect, useCallback, useImperativeHandle, forwardRef } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import CharacterCount from '@tiptap/extension-character-count';
import Underline from '@tiptap/extension-underline';
import { Extension } from '@tiptap/react';
import { Plugin, PluginKey } from '@tiptap/pm/state';
import { Decoration, DecorationSet } from '@tiptap/pm/view';
import { Button, Flex, Tooltip, Typography, theme, Divider } from 'antd';
import {
  BoldOutlined, ItalicOutlined, UnderlineOutlined,
  OrderedListOutlined, UnorderedListOutlined,
  UndoOutlined, RedoOutlined, LineOutlined,
} from '@ant-design/icons';

const { Text } = Typography;

// 错误标记数据
export interface ErrorMark {
  from: number;
  to: number;
  severity: 'error' | 'warning' | 'info';
  message: string;
}

const errorHighlightKey = new PluginKey('errorHighlight');

// Tiptap 扩展：错误高亮装饰
const ErrorHighlight = Extension.create({
  name: 'errorHighlight',

  addOptions() {
    return {
      errorMarks: [] as ErrorMark[],
    };
  },

  addProseMirrorPlugins() {
    const ext = this;
    return [
      new Plugin({
        key: errorHighlightKey,
        props: {
          decorations(state) {
            const marks: ErrorMark[] = ext.options.errorMarks;
            if (!marks || marks.length === 0) return DecorationSet.empty;

            const decorations: Decoration[] = [];
            for (const mark of marks) {
              if (mark.from >= 0 && mark.to <= state.doc.content.size && mark.from < mark.to) {
                const className =
                  mark.severity === 'error' ? 'error-decoration' :
                  mark.severity === 'warning' ? 'warning-decoration' : 'info-decoration';
                decorations.push(
                  Decoration.inline(mark.from, mark.to, {
                    class: className,
                    title: mark.message,
                  })
                );
              }
            }
            return DecorationSet.create(state.doc, decorations);
          },
        },
      }),
    ];
  },
});

export interface RichEditorRef {
  getHTML: () => string;
  getText: () => string;
  setContent: (content: string) => void;
  setErrorMarks: (marks: ErrorMark[]) => void;
  clearErrorMarks: () => void;
  getEditor: () => ReturnType<typeof useEditor>;
}

interface RichEditorProps {
  content?: string;
  placeholder?: string;
  disabled?: boolean;
  onChange?: (text: string) => void;
  minHeight?: number;
  showToolbar?: boolean;
  errorMarks?: ErrorMark[];
}

const RichEditor = forwardRef<RichEditorRef, RichEditorProps>(
  ({ content = '', placeholder, disabled = false, onChange, minHeight = 400, showToolbar = true, errorMarks = [] }, ref) => {
    const { token } = theme.useToken();

    const editor = useEditor({
      extensions: [
        StarterKit.configure({
          heading: { levels: [2, 3] },
        }),
        Placeholder.configure({
          placeholder: placeholder || '在这里开始写作...',
        }),
        CharacterCount,
        Underline,
        ErrorHighlight.configure({
          errorMarks,
        }),
      ],
      content: content || '',
      editable: !disabled,
      onUpdate: ({ editor: e }) => {
        onChange?.(e.getText());
      },
    });

    // 同步 errorMarks 更新
    useEffect(() => {
      if (editor) {
        // 更新扩展的 options 并触发重新渲染
        const ext = editor.extensionManager.extensions.find(e => e.name === 'errorHighlight');
        if (ext) {
          ext.options.errorMarks = errorMarks;
          // 强制 ProseMirror 重新计算 decorations
          editor.view.dispatch(editor.state.tr);
        }
      }
    }, [editor, errorMarks]);

    useEffect(() => {
      if (editor && disabled !== undefined) {
        editor.setEditable(!disabled);
      }
    }, [editor, disabled]);

    // Sync external content changes
    useEffect(() => {
      if (editor && content !== undefined) {
        const currentText = editor.getText();
        if (content !== currentText) {
          editor.commands.setContent(content || '');
        }
      }
    }, [content, editor]);

    const setContent = useCallback((newContent: string) => {
      if (editor) {
        editor.commands.setContent(newContent || '');
      }
    }, [editor]);

    const setErrorMarksInternal = useCallback((marks: ErrorMark[]) => {
      if (editor) {
        const ext = editor.extensionManager.extensions.find(e => e.name === 'errorHighlight');
        if (ext) {
          ext.options.errorMarks = marks;
          editor.view.dispatch(editor.state.tr);
        }
      }
    }, [editor]);

    const clearErrorMarks = useCallback(() => {
      setErrorMarksInternal([]);
    }, [setErrorMarksInternal]);

    useImperativeHandle(ref, () => ({
      getHTML: () => editor?.getHTML() || '',
      getText: () => editor?.getText() || '',
      setContent,
      setErrorMarks: setErrorMarksInternal,
      clearErrorMarks,
      getEditor: () => editor,
    }));

    if (!editor) return null;

    const charCount = editor.storage.characterCount?.characters() || 0;

    return (
      <div
        style={{
          border: `1px solid ${token.colorBorder}`,
          borderRadius: token.borderRadius,
          overflow: 'hidden',
          background: disabled ? token.colorBgLayout : token.colorBgContainer,
        }}
      >
        {showToolbar && !disabled && (
          <Flex
            align="center"
            gap={2}
            style={{
              padding: '4px 8px',
              borderBottom: `1px solid ${token.colorBorderSecondary}`,
              background: token.colorBgLayout,
              flexWrap: 'wrap',
            }}
          >
            <ToolBtn
              icon={<BoldOutlined />}
              title="加粗"
              active={editor.isActive('bold')}
              onClick={() => editor.chain().focus().toggleBold().run()}
            />
            <ToolBtn
              icon={<ItalicOutlined />}
              title="斜体"
              active={editor.isActive('italic')}
              onClick={() => editor.chain().focus().toggleItalic().run()}
            />
            <ToolBtn
              icon={<UnderlineOutlined />}
              title="下划线"
              active={editor.isActive('underline')}
              onClick={() => editor.chain().focus().toggleUnderline().run()}
            />
            <Divider type="vertical" style={{ margin: '0 4px' }} />
            <ToolBtn
              icon={<UnorderedListOutlined />}
              title="无序列表"
              active={editor.isActive('bulletList')}
              onClick={() => editor.chain().focus().toggleBulletList().run()}
            />
            <ToolBtn
              icon={<OrderedListOutlined />}
              title="有序列表"
              active={editor.isActive('orderedList')}
              onClick={() => editor.chain().focus().toggleOrderedList().run()}
            />
            <ToolBtn
              icon={<LineOutlined />}
              title="分隔线"
              onClick={() => editor.chain().focus().setHorizontalRule().run()}
            />
            <Divider type="vertical" style={{ margin: '0 4px' }} />
            <ToolBtn
              icon={<UndoOutlined />}
              title="撤销"
              onClick={() => editor.chain().focus().undo().run()}
              disabled={!editor.can().undo()}
            />
            <ToolBtn
              icon={<RedoOutlined />}
              title="重做"
              onClick={() => editor.chain().focus().redo().run()}
              disabled={!editor.can().redo()}
            />

            <div style={{ marginLeft: 'auto' }}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                {charCount.toLocaleString()} 字
              </Text>
            </div>
          </Flex>
        )}

        <EditorContent
          editor={editor}
          style={{
            minHeight,
            padding: 16,
            fontSize: 15,
            lineHeight: 1.8,
            fontFamily: 'Georgia, "Times New Roman", serif',
            cursor: disabled ? 'not-allowed' : 'text',
          }}
        />
      </div>
    );
  }
);

RichEditor.displayName = 'RichEditor';

function ToolBtn({
  icon,
  title,
  active = false,
  disabled = false,
  onClick,
}: {
  icon: React.ReactNode;
  title: string;
  active?: boolean;
  disabled?: boolean;
  onClick: () => void;
}) {
  const { token } = theme.useToken();
  return (
    <Tooltip title={title}>
      <Button
        type="text"
        size="small"
        icon={icon}
        onClick={onClick}
        disabled={disabled}
        style={{
          background: active ? token.colorPrimaryBg : 'transparent',
          color: active ? token.colorPrimary : undefined,
        }}
      />
    </Tooltip>
  );
}

export default RichEditor;
